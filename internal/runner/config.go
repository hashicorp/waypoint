package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// watchConfig sits in a goroutine receiving the new configurations from the
// server.
func (r *Runner) watchConfig(ch <-chan *pb.RunnerConfig) {
	for config := range ch {
		r.handleConfig(config)
	}
}

// handleConfig handles the changes for a single config.
//
// This is NOT thread-safe, but it is safe to handle a configuration
// change in parallel to any other operation.
func (r *Runner) handleConfig(c *pb.RunnerConfig) {
	old := r.config
	r.config = c

	// Store our original environment as a set of config vars. This will
	// let us replace any of these later if the runtime config gets unset.
	if r.originalEnv == nil {
		r.originalEnv = []*pb.ConfigVar{}
		for _, str := range os.Environ() {
			idx := strings.Index(str, "=")
			if idx == -1 {
				continue
			}

			r.originalEnv = append(r.originalEnv, &pb.ConfigVar{
				Name:  str[:idx],
				Value: &pb.ConfigVar_Static{Static: str[idx+1:]},
			})
		}
	}

	// Handle config var changes
	{
		// Setup our original env. This will ensure that we replace the
		// variable if it becomes unset.
		env := map[string]string{}
		for _, v := range r.originalEnv {
			env[v.Name] = v.Value.(*pb.ConfigVar_Static).Static
		}

		if old != nil {
			// Unset any previous config variables. We check if its in env
			// already because if it is, it is an original value and we accept
			// that. This lets unset runtime config get reset back to the
			// original process start env.
			for _, v := range old.ConfigVars {
				if _, ok := env[v.Name]; !ok {
					env[v.Name] = ""
				}
			}
		}

		// Set the config variables
		for _, v := range c.ConfigVars {
			static, ok := v.Value.(*pb.ConfigVar_Static)
			if !ok {
				r.logger.Warn("unknown value type for config var, ignoring",
					"type", fmt.Sprintf("%T", v.Value))
				continue
			}

			env[v.Name] = static.Static
		}

		// Set them all
		for k, v := range env {
			// We ignore current value so that the log doesn't look messy
			if os.Getenv(k) == v {
				continue
			}

			// Unset if empty
			if v == "" {
				r.logger.Info("unsetting env var", "key", k)
				if err := os.Unsetenv(k); err != nil {
					r.logger.Warn("error unsetting config var", "key", k, "err", err)
				}

				continue
			}

			// Set
			r.logger.Info("setting env var", "key", k)
			if err := os.Setenv(k, v); err != nil {
				r.logger.Warn("error setting config var", "key", k, "err", err)
			}
		}
	}
}

func (r *Runner) recvConfig(
	ctx context.Context,
	client pb.Waypoint_RunnerConfigClient,
	ch chan<- *pb.RunnerConfig,
) {
	log := r.logger.Named("config_recv")
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
			if err == io.EOF || err == context.Canceled {
				return
			}

			log.Error("error receiving configuration, exiting", "err", err)
			return
		}

		log.Info("new configuration received")
		ch <- resp.Config
	}
}
