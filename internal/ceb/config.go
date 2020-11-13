package ceb

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"sort"

	"github.com/hashicorp/go-argmapper"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	sdkpb "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (ceb *CEB) initConfigStream(ctx context.Context, cfg *config) error {
	log := ceb.logger.Named("config")

	// Start the watcher. This will do nothing until anything is sent on the
	// channel so we can start it early. We share the same channel across
	// config reconnects.
	ch := make(chan *pb.EntrypointConfig)
	go ceb.watchConfig(ctx, log, cfg, ch)

	// Start the config receiver. This will connect ot the EntrypointConfig
	// endpoint and start receiving data. This will reconnect on failure.
	go ceb.initConfigStreamReceiver(ctx, log, cfg, ch, false)

	return nil
}

func (ceb *CEB) initConfigStreamReceiver(
	ctx context.Context,
	log hclog.Logger,
	cfg *config,
	ch chan<- *pb.EntrypointConfig,
	isRetry bool,
) error {
	// On retry we always mark the child process ready so we can begin executing
	// any staged child command. We don't do this on non-retries because we
	// still have hope that we can talk to the server and get our initial config.
	if isRetry {
		ceb.markChildCmdReady()
	}

	// wait for initial server connection
	serverClient := ceb.waitClient()
	if serverClient == nil {
		return ctx.Err()
	}

	// Open our log stream
	log.Debug("registering instance, requesting config")
	client, err := serverClient.EntrypointConfig(ctx, &pb.EntrypointConfigRequest{
		DeploymentId: ceb.deploymentId,
		InstanceId:   ceb.id,
	}, grpc.WaitForReady(isRetry || cfg.ServerRequired))
	if err != nil {
		// If the server is unavailable and this is our first time, then
		// we just start this up in the background in retry mode and allow
		// the startup to continue so we don't block the child process starting.
		if status.Code(err) == codes.Unavailable {
			log.Error("error connecting to Waypoint server, will retry but startup " +
				"child command without initial settings")
			go ceb.initConfigStreamReceiver(ctx, log, cfg, ch, true)
			return nil
		}

		return err
	}

	// We never send anything
	client.CloseSend()

	// Start the goroutine that waits for all other configs
	go ceb.recvConfig(ctx, client, ch, func() error {
		return ceb.initConfigStreamReceiver(ctx, log, cfg, ch, true)
	})

	return nil
}

// watchConfig sits in a goroutine receiving the new configurations from the
// server.
func (ceb *CEB) watchConfig(
	ctx context.Context,
	log hclog.Logger,
	cfg *config,
	ch <-chan *pb.EntrypointConfig,
) {
	log = log.Named("watcher")

	// Keep track of our currently executing command information so that
	// we can diff properly to determine if we need to restart.
	currentCmd := ceb.copyCmd(ceb.childCmdBase)

	// We only init the URL service once. In the future, we can do diffing
	// and support automatically reinitializing if the URL service changes.
	didInitURL := false

	for {
		select {
		case <-ctx.Done():
			log.Warn("exiting, context ended")
			return

		case config := <-ch:
			if !didInitURL {
				didInitURL = true

				// If we have URL service configuration, start it. We start this in a goroutine
				// since we don't need to block starting up our application on this.
				if url := config.UrlService; url != nil {
					go func() {
						if err := ceb.initURLService(ctx, cfg.URLServicePort, url); err != nil {
							log.Warn("error starting URL service", "err", err)
						}
					}()
				} else {
					log.Debug("no URL service configuration, will not register with URL service")
				}
			}

			// Start the exec sessions if we have any
			if len(config.Exec) > 0 {
				ceb.startExecGroup(config.Exec)
			}

			// Configure our env vars for the child command.
			ceb.handleChildCmdConfig(log, config, currentCmd)
			ceb.markChildCmdReady()
		}
	}
}

func (ceb *CEB) tickChildCmdConfig(
	log hclog.Logger,
	last *exec.Cmd,
	vars []*pb.ConfigVar,
) {
	// Build up our env vars. We append to our base command. We purposely
	// make a capacity of our _last_ command to try to avoid allocations
	// in the common case (same env).
	base := ceb.childCmdBase
	env := make([]string, len(base.Env), len(last.Env))
	copy(env, base.Env)

	// We want to accumulate the static vars directly on the env and then
	// store the dynamic ones in a mapping of source to vars so we can more
	// easily process those.
	dynamic := map[string][]*component.ConfigRequest{}
	for _, cv := range vars {
		switch v := cv.Value.(type) {
		case *pb.ConfigVar_Static:
			env = append(env, cv.Name+"="+v.Static)

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

	// Determine if there are any config changes and mark which are changed.
	changed := map[string]struct{}{}
	// TODO

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

	// Go through each and read our configurations. Note that ConfigSourcers
	// are documented to note that Read will be called frequently so caching
	// is expected within the sourcer itself.
	for k, reqs := range dynamic {
		L := log.With("source", k)
		s := ceb.configPlugins[k].Component.(component.ConfigSourcer)

		// If the configuration has changed for this plugin, we call Stop.
		if _, ok := changed[k]; ok {
			_, err := ceb.callDynamicFunc(L, s.StopFunc())
			if err != nil {
				// We just continue on error but warn the user. We continue
				// because stop really shouldn't do much here on the plugin
				// side except maybe clear some caches, so errors are unlikely.
				L.Warn("error stopping config source", "err", err)
			}
		}

		// Next, call Read
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

	// Sort the env vars we have so that we can compare reliably
	sort.Strings(env)

	// If the env vars have not changed, we haven't changed. We do this
	// using basic DeepEqual since we always sort the strings here.
	if reflect.DeepEqual(last.Env, env) {
		return
	}

	log.Info("env vars changed, sending new child command")

	// Update the env vars
	last.Env = env

	// Send the new command
	ceb.childCmdCh <- ceb.copyCmd(last)
}

func (ceb *CEB) handleChildCmdConfig(
	log hclog.Logger,
	config *pb.EntrypointConfig,
	last *exec.Cmd,
) {
	// Build up our env vars. We append to our base command. We purposely
	// make a capacity of our _last_ command to try to avoid allocations
	// in the command case (same env).
	base := ceb.childCmdBase
	env := make([]string, len(base.Env), len(last.Env))
	copy(env, base.Env)
	for _, cv := range config.EnvVars {
		static, ok := cv.Value.(*pb.ConfigVar_Static)
		if !ok {
			log.Warn("unknown config value type received, ignoring",
				"type", fmt.Sprintf("%T", cv.Value))
			continue
		}

		env = append(env, cv.Name+"="+static.Static)
	}
	sort.Strings(env)

	// If the env vars have not changed, we haven't changed. We do this
	// using basic DeepEqual since we always sort the strings here.
	if reflect.DeepEqual(last.Env, env) {
		return
	}

	log.Info("env vars changed, sending new child command")

	// Update the env vars
	last.Env = env

	// Send the new command
	ceb.childCmdCh <- ceb.copyCmd(last)
}

func (ceb *CEB) recvConfig(
	ctx context.Context,
	client pb.Waypoint_EntrypointConfigClient,
	ch chan<- *pb.EntrypointConfig,
	reconnect func() error,
) {
	log := ceb.logger.Named("config_recv")
	defer log.Trace("exiting receive goroutine")

	for {
		// If the context is closed, exit
		if ctx.Err() != nil {
			return
		}

		// Wait for the next configuration
		resp, err := client.Recv()
		if err != nil {
			// If we get the unavailable error then the connection died.
			// We restablish the connection.
			if status.Code(err) == codes.Unavailable {
				log.Error("ceb disconnected from server, attempting reconnect")
				err = reconnect()

				// If we successfully reconnected, then exit this.
				if err == nil {
					return
				}
			}

			if err != nil {
				log.Error("error receiving configuration, exiting", "err", err)
				return
			}
		}

		log.Info("new configuration received")
		ch <- resp.Config
	}
}
