package ceb

import (
	"context"
	"fmt"
	"os/exec"
	"reflect"
	"sort"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (ceb *CEB) initConfigStream(ctx context.Context, cfg *config, isRetry bool) error {
	log := ceb.logger.Named("config")

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
			go ceb.initConfigStream(ctx, cfg, true)
			return nil
		}

		return err
	}

	// We never send anything
	client.CloseSend()

	// Start the watcher
	ch := make(chan *pb.EntrypointConfig)
	go ceb.watchConfig(log, ch, ctx, cfg)

	// Start the goroutine that waits for all other configs
	go ceb.recvConfig(ctx, client, ch, func() error {
		return ceb.initConfigStream(ctx, cfg, true)
	})

	return nil
}

// watchConfig sits in a goroutine receiving the new configurations from the
// server.
func (ceb *CEB) watchConfig(
	log hclog.Logger,
	ch <-chan *pb.EntrypointConfig,
	ctx context.Context,
	cfg *config,
) {
	// Keep track of our currently executing command information so that
	// we can diff properly to determine if we need to restart.
	currentCmd := ceb.copyCmd(ceb.childCmdBase)

	// We only init the URL service once. In the future, we can do diffing
	// and support automatically reinitializing if the URL service changes.
	didInitURL := false

	for config := range ch {
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
	defer close(ch)

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
