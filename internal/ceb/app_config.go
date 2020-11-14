package ceb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/r3labs/diff"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdkpb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var (
	// appConfigRefreshPeriod is the interval between checking for new
	// config values. In a steady state, configuration NORMALLY doesn't
	// change so this is set fairly high to avoid unnecessary load on
	// dynamic config sources.
	//
	// NOTE(mitchellh): In the future, we'd like to build a way for
	// config sources to edge-trigger when changes happen to prevent
	// this refresh.
	appConfigRefreshPeriod = 15 * time.Second
)

func (ceb *CEB) watchAppConfig(
	ctx context.Context,
	log hclog.Logger,
	inCh <-chan []*pb.ConfigVar,
	outCh chan<- []string,
) {
	// prevVars keeps track of the previous seen variables sent on inCh.
	// We do some diffing to prevent unnecessary config fetching or command
	// restarting and this is how we account for that.
	var prevVars []*pb.ConfigVar
	var prevVarsChanged map[string]bool

	// prevEnv keeps track of the last set of env vars we computed. We do
	// this to compare and prevent unnecessarilly restarting the command.
	var prevEnv []string

	// static keeps track of the static env vars that we have and dynamic
	// keeps track of all the dynamic configurations that we have.
	var static []string
	var dynamic map[string][]*component.ConfigRequest

	// refreshCh will be sent a message when we want to refresh our
	// configuration. We default to nil so that we do nothing until
	// we receive our first set of variables (the <-inCh case below).
	var refreshCh <-chan time.Time

	// prevSent is flipped to true once we send our first set of compiled
	// env vars to the outCh. We have to keep track of this because there is
	// an expectation that we will always send an initial set of configs.
	prevSent := false

	for {
		select {
		// Case: context is over, we're done
		case <-ctx.Done():
			return

		// Case: caller sends us a new set of variables
		case newVars := <-inCh:
			// If the variables are the same as the last set, then we do nothing.
			if prevSent && ceb.sameAppConfig(log, prevVars, newVars) {
				log.Trace("got var update but ignoring since they're the same")
				continue
			}

			// New variables, track it and immediately trigger a refresh
			log.Debug("new config variables received, scheduling refresh")
			prevVars = newVars
			refreshChNow := make(chan time.Time)
			close(refreshChNow)
			refreshCh = refreshChNow

			// Split the static and dynamic out here since this is something
			// we're going to need often so we precompute it once.
			dynamicOld := dynamic
			static, dynamic = ceb.splitAppConfig(log, newVars)

			// We need to do a diff of if any dynamic var config changed.
			prevVarsChanged = ceb.diffDynamicAppConfig(log, dynamicOld, dynamic)

		// Case: timer fires to refresh our dynamic variable sources
		case <-refreshCh:
			// Set the refreshCh to nil immediately so we never get in an
			// infinite refresh situation on a closed channel.
			refreshCh = nil

			// Get our new env vars
			log.Trace("refreshing app configuration")
			newEnv := ceb.buildAppConfig(log, static, dynamic, prevVarsChanged)
			sort.Strings(newEnv)

			// Mark that we aren't seeing any new vars anymore. This speeds up
			// future buildAppConfig calls since it prevents all the diff logic
			// from happening to detect what plugins need to call Stop.
			prevVarsChanged = nil

			// Setup our next refresh. This "leaks" timers in the scenario
			// we get a lot of variable changes but that is an unlikely case.
			refreshCh = time.After(appConfigRefreshPeriod)

			// Compare our new env and old env. prevEnv is already sorted.
			if prevSent && reflect.DeepEqual(prevEnv, newEnv) {
				log.Trace("app configuration unchanged")
				continue
			}

			// New env vars!
			log.Debug("new configuration computed, sending to child process manager")
			prevEnv = newEnv
			select {
			case outCh <- newEnv:
				prevSent = true

			case <-ctx.Done():
				return
			}
		}
	}
}

// sameAppConfig returns true if the vars and prevVars represent the
// same application configuration.
func (ceb *CEB) sameAppConfig(
	log hclog.Logger,
	vars []*pb.ConfigVar,
	prevVars []*pb.ConfigVar,
) bool {
	// If the lengths are different we can fast track this whole thing.
	if len(vars) != len(prevVars) {
		return false
	}

	// Start by sorting the variables by name.
	sort.Slice(vars, configVarSortFunc(vars))
	sort.Slice(vars, configVarSortFunc(prevVars))

	// Marshal to JSON and compare their values. This is a lazy way to diff.
	// If there are any marshalilng errors we just log and return false.
	bytes1, err1 := json.Marshal(vars)
	bytes2, err2 := json.Marshal(prevVars)
	if err1 != nil || err2 != nil {
		log.Warn("error marshaling config vars for comparison, shouldn't happen",
			"err1", err1,
			"err2", err2)
		return false
	}

	return bytes.Equal(bytes1, bytes2)
}

func configVarSortFunc(vars []*pb.ConfigVar) func(i, j int) bool {
	return func(i, j int) bool {
		return vars[i].Name < vars[j].Name
	}
}

// splitAppConfig takes a list of config variables as sent on the wire
// and splits them into a set of static env vars (in KEY=VALUE format already),
// and a map of dynamic config requests keyed by plugin type.
func (ceb *CEB) splitAppConfig(
	log hclog.Logger,
	vars []*pb.ConfigVar,
) (static []string, dynamic map[string][]*component.ConfigRequest) {
	// Split out our static and dynamic here.
	dynamic = map[string][]*component.ConfigRequest{}
	for _, cv := range vars {
		switch v := cv.Value.(type) {
		case *pb.ConfigVar_Static:
			static = append(static, cv.Name+"="+v.Static)

		case *pb.ConfigVar_Dynamic:
			from := v.Dynamic.From
			dynamic[from] = append(dynamic[from], &component.ConfigRequest{
				Name:   cv.Name,
				Config: v.Dynamic.Config,
			})

		default:
			log.Warn("unknown config value type received, ignoring",
				"type", fmt.Sprintf("%T", cv.Value))
		}
	}

	return
}

// diffDynamicAppConfig determines what config source plugins had any
// changes occur between them. These need to be known so that Stop
// can be called and the plugin potentially stopped.
//
// The return value are all the plugins with changes, and the bool value
// is true if the plugin process should also be killed.
func (ceb *CEB) diffDynamicAppConfig(
	log hclog.Logger,
	dynamicOld, dynamicNew map[string][]*component.ConfigRequest,
) map[string]bool {
	log.Trace("calculating changes between old and new config")
	changed := map[string]bool{}

	// Anything in the old and not in the new needs to be stopped.
	for k := range dynamicOld {
		if _, ok := dynamicNew[k]; !ok {
			log.Trace("config source longer in use", "source", k)
			changed[k] = true
		}
	}

	// Go through new. Anything in new and not in old is a change. If
	// it is in both, we have to do a comparison by requests.
	for k := range dynamicNew {
		if _, ok := dynamicOld[k]; !ok {
			log.Trace("config source is new", "source", k)
			changed[k] = false
			continue
		}

		reqsOld := map[string]*component.ConfigRequest{}
		for _, req := range dynamicOld[k] {
			reqsOld[req.Name] = req
		}

		reqsNew := map[string]*component.ConfigRequest{}
		for _, req := range dynamicNew[k] {
			reqsNew[req.Name] = req
		}

		changes, _ := diff.Diff(reqsOld, reqsNew)
		if len(changes) > 0 {
			log.Trace("config source changed", "source", k)
			changed[k] = false
		}
	}

	return changed
}

// buildAppConfig takes the static and dynamic variables and builds up the
// full list of actual env variable values.
func (ceb *CEB) buildAppConfig(
	log hclog.Logger,
	static []string,
	dynamic map[string][]*component.ConfigRequest,
	changed map[string]bool,
) []string {
	// For each dynamic config, we need to launch that plugin if we
	// haven't already.
	for k, _ := range dynamic {
		if _, ok := ceb.configPlugins[k]; ok {
			continue
		}

		// NOTE(mitchellh): For the initial version, we hardcode all our
		// config sourcers directly so there is no actual plugin loading
		// happening. Instead, we're just validating that the plugin is known.
		// In the future, this is roughly where we should hook up plugin loading.
		log.Warn("unknown config source plugin requested", "name", k)
	}

	// Go through the changed plugins first and call Stop.
	for k, kill := range changed {
		raw, ok := ceb.configPlugins[k]
		if !ok {
			continue
		}

		L := log.With("source", k)
		L.Debug("config variables changed, calling Stop")
		s := raw.Component.(component.ConfigSourcer)
		_, err := ceb.callDynamicFunc(L, s.StopFunc())
		if err != nil {
			// We just continue on error but warn the user. We continue
			// because stop really shouldn't do much here on the plugin
			// side except maybe clear some caches, so errors are unlikely.
			L.Warn("error stopping config source", "err", err)
		}

		if kill {
			L.Debug("config variables no longer using this source, killing")

			// End it
			if raw.Close != nil {
				raw.Close()
			}

			// Delete it from our plugins map
			// NOTE(mitchellh): we don't do this right now because we don't
			// actually load plugins yet.
		}
	}

	// If we have no dynamic values, then we just return the static ones.
	if len(dynamic) == 0 {
		return static
	}

	// Ininitialize our result with the static values
	env := make([]string, len(static), len(static)*2)
	copy(env, static)

	// Go through each and read our configurations. Note that ConfigSourcers
	// are documented to note that Read will be called frequently so caching
	// is expected within the sourcer itself.
	for k, reqs := range dynamic {
		L := log.With("source", k)
		s := ceb.configPlugins[k].Component.(component.ConfigSourcer)

		// Next, call Read
		if L.IsTrace() {
			var keys []string
			for _, req := range reqs {
				keys = append(keys, req.Name)
			}
			L.Trace("reading values for keys", "keys", keys)
		}
		result, err := ceb.callDynamicFunc(L, s.ReadFunc(),
			argmapper.Typed(reqs),
		)
		if err != nil {
			L.Warn("error reading configuration values, all will be dropped", "err", err)
			continue
		}

		// Get the result
		if result.Len() != 1 {
			L.Warn("config source should've returned one result, dropping results", "got", result.Len())
			continue
		}
		values, ok := result.Out(0).([]*sdkpb.ConfigSource_Value)
		if !ok {
			L.Warn("config should returned invalid type, dropping",
				"got", fmt.Sprintf("%T", result.Out(0)))
			continue
		}

		// Build a map so that we only include values we care about.
		valueMap := map[string]*sdkpb.ConfigSource_Value{}
		for _, v := range values {
			valueMap[v.Name] = v
		}
		for _, req := range reqs {
			value, ok := valueMap[req.Name]
			if !ok {
				L.Warn("config source didn't populate expected value", "key", req.Name)
				continue
			}

			switch r := value.Result.(type) {
			case *sdkpb.ConfigSource_Value_Value:
				env = append(env, req.Name+"="+r.Value)

			case *sdkpb.ConfigSource_Value_Error:
				st := status.FromProto(r.Error)
				L.Warn("error retrieving config value",
					"key", req.Name,
					"err", st.Err().Error())

			default:
				L.Warn("config value had unknown result type, ignoring",
					"key", req.Name,
					"type", fmt.Sprintf("%T", value.Result))
			}
		}
	}

	return env
}
