package appconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/r3labs/diff"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdkpb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint/internal/pkg/condctx"
	"github.com/hashicorp/waypoint/internal/plugin"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var (
	// defaultRefreshInterval is picked to be long enough to not overstrain
	// systems but short enough that config changes propagate reasonably.
	defaultRefreshInterval = 15 * time.Second
)

// Watcher reads application configuration values and watches for any changes.
//
// The values that the watcher is watching can be added, removed, or updated
// along with any configuration sources (how to read from external systems
// such as Vault).
type Watcher struct {
	log hclog.Logger

	// dynamicEnabled determines whether we allow dynamic sources or not.
	// If this is false, then we ignore all dynamic configs.
	dynamicEnabled bool

	// refreshInterval is the interval between checking for new
	// config values. In a steady state, configuration NORMALLY doesn't
	// change so this is set fairly high to avoid unnecessary load on
	// dynamic config sources.
	//
	// NOTE(mitchellh): In the future, we'd like to build a way for
	// config sources to edge-trigger when changes happen to prevent
	// this refresh.
	refreshInterval time.Duration

	// plugins is a set of plugins that are already launched for
	// config sourcing.
	plugins map[string]*plugin.Instance

	// inSourceCh and inVarCh are the channels that are used to send
	// updated sets of configuration sources and variables to the watch loop.
	inSourceCh chan []*pb.ConfigSource
	inVarCh    chan []*pb.ConfigVar

	// currentCond is used to lock and notify updates for currentEnv.
	currentCond *sync.Cond

	// currentEnv is the list of current environment variable values for
	// the configuration.
	currentEnv []string

	// currentGen is the current "generation" of configuration values. This
	// is incremented by one each time the current config value (currentEnv)
	// are updated. This can be used along with currentCond to detect
	// changes in currentEnv.
	currentGen uint64

	// bgCtx, bgCancel, and bgWg are all used for lifecycle management of
	// background goroutines managed by the watcher. bgCtx can be used to
	// cancel them (via bgCancel), and bgWg can be waited on to ensure
	// everything is stopped.
	bgCtx    context.Context
	bgCancel context.CancelFunc
	bgWg     *sync.WaitGroup
}

// NewWatcher creates a new Watcher instance.
//
// This will immediately start the background goroutine for reading and
// updating configuration values, even if no initial values are provided.
// You must call Close to properly clean up resources used by the Watcher.
func NewWatcher(opts ...Option) (*Watcher, error) {
	var bgWg sync.WaitGroup
	bgCtx, bgCancel := context.WithCancel(context.Background())

	// If we return due to an error, cancel the background context.
	// This won't do anything on success cause we nil out bgCancel.
	defer func() {
		if bgCancel != nil {
			bgCancel()
		}
	}()

	// Build our initial watcher
	w := &Watcher{
		log:             hclog.L(),
		dynamicEnabled:  true,
		refreshInterval: defaultRefreshInterval,
		plugins:         map[string]*plugin.Instance{},
		inSourceCh:      make(chan []*pb.ConfigSource),
		inVarCh:         make(chan []*pb.ConfigVar),
		currentCond:     sync.NewCond(&sync.Mutex{}),
		bgCtx:           bgCtx,
		bgCancel:        bgCancel,
		bgWg:            &bgWg,
	}

	// Use the option pattern to update any options.
	for _, opt := range opts {
		if err := opt(w); err != nil {
			return nil, err
		}
	}

	// Start our background goroutine
	w.bgWg.Add(1)
	go w.watcher(
		bgCtx,
		w.log.Named("watchloop"),
	)

	// Everything is good, nil out bgCancel so our defer doesn't stop us
	bgCancel = nil

	return w, nil
}

// Close stops all the background goroutines that this watcher started.
// This will block until all the background tasks have exited.
func (w *Watcher) Close() error {
	w.bgCancel()
	w.bgWg.Wait()
	return nil
}

// Next returns the next values for the configuration AFTER the given
// iterator value iter. A value of 0 can be used for iter for a first read.
//
// The return value will be the configuration values in env format (KEY=VALUE),
// the current iterator value that you should use with the next call to Next,
// and any error if it occurred.
//
// The ctx parameter can be used for timeouts, cancellation, etc. If the context
// is closed, this will return the context error.
func (w *Watcher) Next(ctx context.Context, iter uint64) ([]string, uint64, error) {
	var cancelFunc func()

	w.currentCond.L.Lock()
	defer w.currentCond.L.Unlock()

	// Wait on the condition var as long as we have the same iterator
	// and the context isn't yet cancelled.
	for w.currentGen == iter && ctx.Err() == nil {
		// If we're waiting, then we want to start a goroutine to notify
		// us if the context closes. We have to do this in a goroutine because
		// cond vars have no other way to wait on a context.
		//
		// We do this in the for loop so that on the fast path where we
		// have an older generation, we just return the value immediately
		// without all the goroutine ceremony.
		if cancelFunc == nil {
			cancelFunc = condctx.Notify(ctx, w.currentCond)
			defer cancelFunc()
		}

		w.currentCond.Wait()
	}

	// If we exited due to context being canceled, exit now.
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	return w.currentEnv, w.currentGen, nil
}

// UpdateSources updates the configuration sources for the watcher. The
// behavior and semantics are identical to UpdateVars but for configuration
// sources, so please see the documentation for UpdateVars for more details.
func (w *Watcher) UpdateSources(ctx context.Context, v []*pb.ConfigSource) error {
	select {
	case w.inSourceCh <- v:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

// UpdateVars updates the variables for the watcher. This replaces all
// the previous set variables.
//
// This may block for some time waiting for the update loop to accept
// our changes. The ctx parameter can be used as a timeout. If the context
// is cancelled, the error returned will be the context error.
func (w *Watcher) UpdateVars(ctx context.Context, v []*pb.ConfigVar) error {
	select {
	case w.inVarCh <- v:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Watcher) notify(
	ctx context.Context,
	ch chan<- []string,
) {
	// lastGen is the last generation we saw. We always set this to zero
	// so we get an initial value sent (first value is 1).
	var lastGen uint64 = 0

	for {
		newEnv, nextGen, err := w.Next(ctx, lastGen)
		if err != nil {
			// This case covers context cancellation as well since
			// Next returns the context error on cancellation.
			return
		}

		lastGen = nextGen
		select {
		case ch <- newEnv:
			// Sent successfully

		case <-ctx.Done():
			// Context over, return
			return
		}
	}
}

// watcher is the main watch loop that waits for changes in configuration
// or configuration sources and sends the resulting set of environment variables
// on the output channel.
//
// Callers must always add one to w.bgWg prior to calling this.
func (w *Watcher) watcher(
	ctx context.Context,
	log hclog.Logger,
) {
	defer w.bgWg.Done()

	// prevVars keeps track of the previous seen variables sent on inVarCh.
	// We do some diffing to prevent unnecessary config fetching or command
	// restarting and this is how we account for that.
	var prevVars []*pb.ConfigVar
	prevVarsChanged := map[string]bool{}

	// prevEnv keeps track of the last set of env vars we computed. We do
	// this to compare and prevent unnecessarilly restarting the command.
	var prevEnv []string

	// static keeps track of the static env vars that we have and dynamic
	// keeps track of all the dynamic configurations that we have.
	var static []string
	var dynamic map[string][]*component.ConfigRequest
	var dynamicSources map[string]*pb.ConfigSource

	// refreshCh will be sent a message when we want to refresh our
	// configuration. We default to nil so that we do nothing until
	// we receive our first set of variables (the <-inVarCh case below).
	//
	// coalesceCh is used when we want to refresh, but allow some time
	// for coalescing of the source/variable channels to occur.
	var refreshCh, coalesceCh <-chan time.Time
	refreshTick := func() {
		// If we haven't scheduled a forced refresh, then schedule that.
		// We will refresh NO MATTER WHAT on this timer and prevents a
		// flurry of config updates from preventing variable refresh.
		if refreshCh == nil {
			refreshCh = time.After(5 * time.Second)
		}

		// Reset our coalesce channel. Using "time.After" here "leaks"
		// timers if we're calling this enough but they're a bunch of timers
		// that reset relatively quickly so let's just let it happen for now.
		coalesceCh = time.After(500 * time.Millisecond)
	}

	// refreshNowCh is just a closed time channel that will trigger
	// a receive immediately. This can be assigned to coalesce or refresh
	// channels to trigger them.
	refreshNowCh := make(chan time.Time)
	close(refreshNowCh)

	// prevSent is flipped to true once we update our first set of compiled
	// env vars to the currentEnv. We have to keep track of this because there is
	// an expectation that we will always set an initial set of configs.
	prevSent := false

	for {
		select {
		// Case: context is over, we're done
		case <-ctx.Done():
			return

		// Case: caller sends us a new set of config source settings
		case newSources := <-w.inSourceCh:
			// Our first pass here is a quick high-level pass to determine if
			// anything is possibly different at all. If it isn't, we just
			// continue on.
			set := map[string]struct{}{}
			diff := map[string]*pb.ConfigSource{}
			for _, source := range newSources {
				set[source.Type] = struct{}{}
				prev, ok := dynamicSources[source.Type]

				// If we haven't seen this before ever, there is a diff.
				// If we have seen this before but the configurations are
				// different then there is also a diff.
				if !ok || prev.Hash != source.Hash {
					diff[source.Type] = source
					continue
				}
			}
			for k := range dynamicSources {
				// Detect if we _removed_ any configurations.
				if _, ok := set[k]; !ok {
					diff[k] = nil
				}
			}
			if len(diff) == 0 {
				log.Trace("got source config update but ignoring since there is no diff")
				continue
			}

			// We have a difference, we now go through and more carefully
			// determine if the difference matters. By "matters" we mean:
			// does it impact dynamic variables we have already fetched? If not,
			// then we just store the config cause when we first fetch we'll
			// grab em. If it does, we have to notify and schedule a refresh
			// because we need to stop and refetch.
			dynamicSources = map[string]*pb.ConfigSource{}
			for k, source := range diff {
				// If we have variables dependent on this config, then
				// we need to mark this as changed. If we don't, then ignore
				// it.
				if len(dynamic[k]) > 0 {
					log.Trace("change in source config, scheduling refresh", "source", k)
					prevVarsChanged[k] = false
				}

				// Ignore nil sources. A nil source means we removed the
				// configuration. We need that so that the above can detect
				// if we have dynamic vars dependent on that but we don't
				// want to store it.
				if source != nil {
					dynamicSources[k] = source
				}
			}

			// If we have changes, schedule a refresh
			if len(prevVarsChanged) > 0 {
				refreshTick()
			}

		// Case: caller sends us a new set of variables
		case newVars := <-w.inVarCh:
			// If the variables are the same as the last set, then we do nothing.
			if prevSent && w.sameAppConfig(log, prevVars, newVars) {
				log.Trace("got var update but ignoring since they're the same")
				continue
			}

			// New variables, track it and immediately trigger a refresh
			log.Debug("new config variables received, scheduling refresh")
			prevVars = newVars
			refreshTick()

			// Split the static and dynamic out here since this is something
			// we're going to need often so we precompute it once.
			dynamicOld := dynamic
			static, dynamic = splitAppConfig(log, newVars)

			// Handle the case we disable dynamics
			if !w.dynamicEnabled && len(dynamic) > 0 {
				log.Debug("dynamic config vars are disabled, ignoring", "n", len(dynamic))
				dynamic = nil
			}

			// We need to do a diff of if any dynamic var config changed.
			// We loop through the result here and set values to true so
			// that we don't clobber changes that inSourceCh receiving may have
			// set. On refresh, we always reset prevVarsChanged to empty.
			for k, v := range w.diffDynamicAppConfig(log, dynamicOld, dynamic) {
				// If it is false, we override it with whatever v we have.
				if !prevVarsChanged[k] {
					prevVarsChanged[k] = v
				}
			}

		// Case: timer fires after a period of time where we have received
		// no other messages and we can now force a refresh.
		case <-coalesceCh:
			// nil the coalesceCh so it isn't called again (until reset)
			coalesceCh = nil

			// set the refreshCh to a closed channel so it triggers ASAP
			refreshCh = refreshNowCh

		// Case: timer fires to refresh our dynamic variable sources
		case <-refreshCh:
			// Set the refreshCh to nil immediately so we never get in an
			// infinite refresh situation on a closed channel.
			refreshCh = nil

			// Set the coalesceCh to nil since we are processing.
			coalesceCh = nil

			// Get our new env vars
			log.Trace("refreshing app configuration")
			newEnv := buildAppConfig(ctx, log,
				w.plugins, static, dynamic, dynamicSources, prevVarsChanged)

			sort.Strings(newEnv)

			// Mark that we aren't seeing any new vars anymore. This speeds up
			// future buildAppConfig calls since it prevents all the diff logic
			// from happening to detect what plugins need to call Stop.
			prevVarsChanged = map[string]bool{}

			// Setup our next refresh. This "leaks" timers in the scenario
			// we get a lot of variable changes but that is an unlikely case.
			refreshCh = time.After(w.refreshInterval)

			// Compare our new env and old env. prevEnv is already sorted.
			if prevSent && reflect.DeepEqual(prevEnv, newEnv) {
				log.Trace("app configuration unchanged")
				continue
			}

			// New env vars!
			log.Debug("new configuration computed")
			prevEnv = newEnv

			// Update our currentEnv
			w.currentCond.L.Lock()
			w.currentEnv = newEnv
			w.currentGen++
			w.currentCond.Broadcast()
			w.currentCond.L.Unlock()

			// We've sent now
			prevSent = true
		}
	}
}

// sameAppConfig returns true if the vars and prevVars represent the
// same application configuration.
func (w *Watcher) sameAppConfig(
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
func splitAppConfig(
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
func (w *Watcher) diffDynamicAppConfig(
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
func buildAppConfig(
	ctx context.Context,
	log hclog.Logger,
	configPlugins map[string]*plugin.Instance,
	static []string,
	dynamic map[string][]*component.ConfigRequest,
	dynamicSources map[string]*pb.ConfigSource,
	changed map[string]bool,
) []string {
	// For each dynamic config, we need to launch that plugin if we
	// haven't already.
	for k := range dynamic {
		if _, ok := configPlugins[k]; ok {
			continue
		}

		// NOTE(mitchellh): For the initial version, we hardcode all our
		// config sourcers directly so there is no actual plugin loading
		// happening. Instead, we're just validating that the plugin is known.
		// In the future, this is roughly where we should hook up plugin loading.
		log.Warn("unknown config source plugin requested", "name", k)
	}

	// erroredSources keeps track of sources that had errors during configuration.
	// If a source is here, we won't load any configs for it.
	erroredSources := map[string]struct{}{}

	// Go through the changed plugins first and call Stop.
	for k, kill := range changed {
		raw, ok := configPlugins[k]
		if !ok {
			continue
		}

		L := log.With("source", k)
		L.Debug("config variables changed, calling Stop")
		s := raw.Component.(component.ConfigSourcer)
		_, err := plugin.CallDynamicFunc(L, s.StopFunc(),
			argmapper.Typed(ctx),
		)
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
			continue
		}

		// Configure the plugin if we have configuration
		configBody := hcl.EmptyBody()
		if s, ok := dynamicSources[k]; ok {
			// We create an hcl.Body by converting the config to JSON first
			// and then using the hcl JSON format. This should always work
			// because our input is a simple map[string]string.
			jsonBytes, err := json.Marshal(s.Config)
			if err != nil {
				panic(err)
			}

			file, diag := hcljson.Parse(jsonBytes, "<config>")
			if diag.HasErrors() {
				panic(diag.Error())
			}

			configBody = file.Body
		}

		diag := component.Configure(raw.Component, configBody, nil)
		if diag.HasErrors() {
			L.Warn("error configuring config source", "err", diag.Error())
			erroredSources[k] = struct{}{}
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

		if _, ok := erroredSources[k]; ok {
			L.Warn("ignoring variables for this source since configuration failed")
			continue
		}

		s := configPlugins[k].Component.(component.ConfigSourcer)

		// Next, call Read
		if L.IsTrace() {
			var keys []string
			for _, req := range reqs {
				keys = append(keys, req.Name)
			}
			L.Trace("reading values for keys", "keys", keys)
		}
		result, err := plugin.CallDynamicFunc(L, s.ReadFunc(),
			argmapper.Typed(ctx),
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
			L.Warn("config source returned invalid type, dropping",
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
