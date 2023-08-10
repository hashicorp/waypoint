// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/factory"
	"github.com/hashicorp/waypoint/internal/plugin"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/tokenutil"
)

var (
	ErrClosed  = errors.New("runner is closed")
	ErrTimeout = errors.New("runner timed out waiting for a job")
)

const (
	// envLogLevel is the env var to set with the log level. This
	// env var matches the Waypoint CLI on purpose. This can be set on
	// the runner process OR via app config (`waypoint config`).
	envLogLevel = "WAYPOINT_LOG_LEVEL"
)

// Runners in Waypoint execute operations. These can be local (the CLI)
// or they can be remote (triggered by some webhook). In either case, they
// share this same underlying implementation.
//
// To use a runner:
//
//  1. Initialize it with New. This will setup some initial state but
//     will not register with the server or run jobs.
//
//  2. Start the runner with "Start". This will register the runner and
//     kick off some management goroutines. This will not execute any jobs.
//
//  3. Run a single job with "Accept". This is named to be similar to a
//     network listener "accepting" a connection. This will request a single
//     job from the Waypoint server, block until one is available, and execute
//     it. Repeat this call for however many jobs you want to execute.
//
//  4. Clean up with "Close". This will gracefully exit the runner, waiting
//     for any running jobs to finish.
type Runner struct {
	id          string
	logger      hclog.Logger
	client      pb.WaypointClient
	cookie      string
	cleanupFunc func()
	runner      *pb.Runner
	factories   map[component.Type]*factory.Factory
	ui          terminal.UI
	local       bool
	tempDir     string
	stateDir    string

	// protects whether or not the runner is active or not.
	runningCond *sync.Cond
	runningJobs int

	runningCtx    context.Context
	runningCancel func()

	enableDynConfig bool

	// stateCond and its associated locker are used to protect all the
	// state-prefixed fields. These state fields can be watched using this
	// cond for state changes in the runner. Anyone waiting on stateCond should
	// also verify the context didn't cancel. The stateCond will be broadcasted
	// when the root context cancels.
	stateCond       *sync.Cond
	stateConfig     uint64 // config stream is connected, increments for each reconnect
	stateConfigOnce uint64 // >0 once we process config once, success or error
	stateExit       uint64 // >0 when exiting

	// config is the current runner config.
	config      *pb.RunnerConfig
	originalEnv []*pb.ConfigVar

	acceptTimeout time.Duration

	// configPlugins is the mapping of config source type to launched plugin.
	// Note this is not currently configurable and we just statically set
	// this to `plugin.ConfigSourcers`. Everything is set up so that
	// in the future it can be configurable though.
	configPlugins map[string]*plugin.Instance

	// noopCh is used in tests only. This will cause any noop operations
	// to block until this channel is closed.
	noopCh <-chan struct{}
}

// New initializes a new runner.
//
// You must call Start to start the runner and register with the Waypoint
// server. See the Runner struct docs for more details.
func New(opts ...Option) (*Runner, error) {
	// Our default runner
	runner := &Runner{
		logger: hclog.L(),
		factories: map[component.Type]*factory.Factory{
			component.MapperType:         plugin.BaseFactories[component.MapperType],
			component.BuilderType:        plugin.BaseFactories[component.BuilderType],
			component.RegistryType:       plugin.BaseFactories[component.RegistryType],
			component.PlatformType:       plugin.BaseFactories[component.PlatformType],
			component.ReleaseManagerType: plugin.BaseFactories[component.ReleaseManagerType],
			component.TaskLauncherType:   plugin.BaseFactories[component.TaskLauncherType],
		},
		stateCond: sync.NewCond(&sync.Mutex{}),
	}

	runner.runningCond = sync.NewCond(new(sync.Mutex))
	runner.runningCtx, runner.runningCancel = context.WithCancel(context.Background())

	// Setup our default config sourcers.
	runner.configPlugins = plugin.ConfigSourcers

	// Build our config
	var cfg config
	for _, o := range opts {
		err := o(runner, &cfg)
		if err != nil {
			return nil, err
		}
	}

	// If we have a state directory, then load that.
	if dir := runner.stateDir; dir != "" {
		if err := verifyStateDir(runner.logger, dir); err != nil {
			return nil, err
		}
	}

	// If the options didn't populate id, then we do so now.
	if runner.id == "" {
		// If we have an ID in state, use that.
		stateId, err := runner.stateGetId()
		if err != nil {
			return nil, err
		}
		if stateId != "" {
			runner.logger.Info("loaded ID from state", "id", stateId)
			runner.id = stateId
		}

		// Create our ID if we still have no ID
		if runner.id == "" {
			id, err := server.Id()
			if err != nil {
				return nil, status.Errorf(codes.Internal,
					"failed to generate unique ID: %s", err)
			}

			// Persist it
			if err := runner.statePutId(id); err != nil {
				return nil, err
			}

			runner.logger.Info("generated a new runner ID", "id", id)
			runner.id = id
		}
	}

	// If we were given a cookie, configure our context to have the cookie
	// for all API requests.
	if v := runner.cookie; v != "" {
		runner.runningCtx = metadata.NewOutgoingContext(
			runner.runningCtx,
			metadata.New(map[string]string{
				"wpcookie": v,
			}),
		)
	}

	runner.runner = &pb.Runner{
		Id:       runner.id,
		ByIdOnly: cfg.byIdOnly,
		Labels:   cfg.labels,
	}

	// Determine what kind of remote runner we are
	if cfg.odr {
		runner.runner.Kind = &pb.Runner_Odr{
			Odr: &pb.Runner_ODR{
				ProfileId: cfg.odrProfileId,
			},
		}
	} else if runner.local {
		runner.runner.Kind = &pb.Runner_Local_{
			Local: &pb.Runner_Local{},
		}
	} else {
		// If this runner isn't ODR or Local, by process of elimination it must be a "static" remote runner.
		// We don't currently have a method to indicate this explicitly.
		runner.runner.Kind = &pb.Runner_Remote_{
			Remote: &pb.Runner_Remote{},
		}
	}

	// Setup our runner components list
	for t, f := range runner.factories {
		for _, n := range f.Registered() {
			runner.runner.Components = append(runner.runner.Components, &pb.Component{
				Type: pb.Component_Type(t),
				Name: n,
			})
		}
	}

	runner.logger.Debug("Created runner", "id", runner.id)

	return runner, nil
}

// Id returns the runner ID.
func (r *Runner) Id() string {
	return r.id
}

// Start starts the runner by registering the runner with the Waypoint
// server. This will spawn goroutines for management. This will return after
// registration so this should not be executed in a goroutine.
func (r *Runner) Start(ctx context.Context) error {
	if r.readState(&r.stateExit) > 0 {
		return ErrClosed
	}

	log := r.logger

	// We set this to true if we're going through the adoption process.
	adopt := false
	tokenCtx := ctx

	// Check if we have a token in our state directory. If we do, then put
	// it on our context (so we use it) and mark that we want to adopt. It
	// may be counterintuitive that we want to adopt if we have a token,
	// but we use that as a way to verify the token and in the future do rotation.
	if t, err := r.stateGetToken(); err != nil {
		return err
	} else if t != "" {
		adopt = true
		log.Debug("will use prior token from state directory")
		tokenCtx = tokenutil.TokenWithContext(tokenCtx, t)
		r.runningCtx = tokenutil.TokenWithContext(r.runningCtx, t)
	}

	// If we have a cookie set, we always adopt.
	if r.cookie != "" {
		adopt = true
		tokenCtx = metadata.NewOutgoingContext(
			tokenCtx,
			metadata.New(map[string]string{
				"wpcookie": r.cookie,
			}),
		)
	} else if !adopt {
		log.Warn("cookie not set for runner, will skip adoption process")
	}

	if adopt {
		// Register and initialize the adoption flow (if necessary) by requesting
		// our token.
		retry := false
		for {
			log.Debug("requesting token with RunnerToken (initiates adoption)")
			tokenResp, err := r.client.RunnerToken(tokenCtx, &pb.RunnerTokenRequest{
				Runner: r.runner,
			}, grpc.WaitForReady(retry))
			if err != nil {
				if status.Code(err) == codes.Unavailable {
					log.Warn("server down during adoption, will attempt reconnect")
					retry = true
					continue
				}
				return err
			}

			if tokenResp != nil && tokenResp.Token != "" {
				// If we received a token, then we replace our token with that.
				// It is possible that we do NOT have a token, because our current
				// token is already valid.
				log.Debug("runner adoption complete, new token received")
				r.runningCtx = tokenutil.TokenWithContext(r.runningCtx, tokenResp.Token)

				// Persist our token
				if err := r.statePutToken(tokenResp.Token); err != nil {
					return err
				}
			} else {
				log.Debug("runner token is already valid, using same token")
			}

			break
		}
	}

	// Note from here on forward, we purposely switch to runningCtx instead
	// of the parameter ctx because everything here on is async and long-running
	// and runningCtx is tied to the full struct lifecycle rather than this
	// single func call.

	// Start our configuration
	log.Debug("starting RunnerConfig stream")
	if err := r.initConfigStream(r.runningCtx); err != nil {
		return err
	}

	// Wait for initial registration
	log.Debug("waiting for registration")
	if r.waitState(&r.stateConfig) {
		return status.Errorf(codes.Internal, "early exit while waiting for first config")
	}

	// Wait for the initial configuration to be set
	log.Debug("runner registered, waiting for first config processing")
	if r.waitState(&r.stateConfigOnce) {
		return status.Errorf(codes.Internal, "early exit while waiting for first config processing")
	}

	log.Info("runner registered with server and ready")
	return nil
}

// Close gracefully exits the runner. This will wait for any pending
// job executions to complete and then deregister the runner. After
// this is called, Start and Accept will no longer function and will
// return errors immediately.
func (r *Runner) Close() error {
	r.runningCond.L.Lock()
	defer r.runningCond.L.Unlock()

	// Wait for all the jobs to finish before we set the shutdown flag.
	for r.runningJobs > 0 {
		r.runningCond.Wait()
	}

	// Mark we're exiting
	r.incrState(&r.stateExit)

	// Cancel the context that is used by Accept to wait on the RunnerJobStream.
	// This interrupts the Recv() call so that any goroutine running Accept()
	// knows that the Runner is closed now and can exit, avoiding a
	// goroutine leak.
	r.runningCancel()

	// Run any cleanup necessary
	if f := r.cleanupFunc; f != nil {
		f()
	}

	return nil
}

// waitState waits for the given state to be set to an initial value.
func (r *Runner) waitState(state *uint64) bool {
	return r.waitStateGreater(state, 0)
}

// waitStateGreater waits for the given state to increment above the value
// v. This can be used with the fields such as stateConfig to detect
// when a reconnect occurs.
func (r *Runner) waitStateGreater(state *uint64, v uint64) bool {
	r.stateCond.L.Lock()
	defer r.stateCond.L.Unlock()
	for *state <= v && r.stateExit == 0 {
		r.stateCond.Wait()
	}

	return r.stateExit > 0
}

// readState reads the current state value. This can be used with
// waitStateGreater to detect a change in state.
func (r *Runner) readState(state *uint64) uint64 {
	r.stateCond.L.Lock()
	defer r.stateCond.L.Unlock()

	// note: we don't use sync/atomic because the writer can't use
	// sync/atomic (see incrState)
	return *state
}

// incrState increments the value of a state variable. The first time
// this is called will also trigger waitState.
func (r *Runner) incrState(state *uint64) {
	r.stateCond.L.Lock()
	defer r.stateCond.L.Unlock()

	// Note: we don't use sync/atomic because we want to pair the increment
	// with the condition variable broadcast. The broadcast requires a lock
	// anyways so there is no need to bring in atomic ops.
	*state += 1
	r.stateCond.Broadcast()
}

type config struct {
	byIdOnly     bool
	odr          bool
	odrProfileId string
	token        string
	labels       map[string]string
}

type Option func(*Runner, *config) error

// WithClient sets the client directly. In this case, the runner won't
// attempt any connection at all regardless of other configuration (env
// vars or waypoint config file). This will be used.
//
// If this is specified, the client MUST use a tokenutil.ContextToken
// type for the PerRPCCredentials setting. This package and others will use
// context overrides for the token. If you do not use this, things will break.
func WithClient(client pb.WaypointClient) Option {
	return func(r *Runner, cfg *config) error {
		r.client = client
		return nil
	}
}

// WithComponentFactory sets a factory for a component type. If this isn't set for
// a component type, then the builtins will be used.
func WithComponentFactory(t component.Type, f *factory.Factory) Option {
	return func(r *Runner, cfg *config) error {
		r.factories[t] = f
		return nil
	}
}

// WithLogger sets the logger that the runner will use. If this isn't
// set it uses hclog.L().
func WithLogger(logger hclog.Logger) Option {
	return func(r *Runner, cfg *config) error {
		r.logger = logger
		return nil
	}
}

// WithLocal sets the runner to local mode. This only changes the UI
// behavior to use the given UI. If ui is nil then the normal streamed
// UI will be used.
func WithLocal(ui terminal.UI) Option {
	return func(r *Runner, cfg *config) error {
		r.local = true
		r.ui = ui
		return nil
	}
}

// ByIdOnly sets it so that only jobs that target this runner by specific
// ID may be assigned.
func ByIdOnly() Option {
	return func(r *Runner, cfg *config) error {
		cfg.byIdOnly = true
		return nil
	}
}

// WithODR configures this runner to be an on-demand runner. This
// will flag this to the server on registration.
func WithODR(profileId string) Option {
	return func(r *Runner, cfg *config) error {
		cfg.odr = true
		cfg.odrProfileId = profileId
		return nil
	}
}

func WithDynamicConfig(set bool) Option {
	return func(r *Runner, cfg *config) error {
		r.enableDynConfig = set
		return nil
	}
}

// WithId sets the id of the runner directly. This isused when the when the server
// is expecting the runner to use a certain ID, such as when used via ondemand runners.
func WithId(id string) Option {
	return func(r *Runner, cfg *config) error {
		r.id = id
		return nil
	}
}

// WithCookie sets the cookie to send with all API requests. If this cookie
// does not match the remote server, API requests will fail.
//
// A cookie is REQUIRED FOR ADOPTION. If this is not set, the adoption process
// will be skipped and only pre-adoption (a preset token) will work.
func WithCookie(v string) Option {
	return func(r *Runner, cfg *config) error {
		r.cookie = v
		return nil
	}
}

// WithStateDir sets the state directory. This directory is used for runner
// state between restarts. This is optional, a runner can be stateless, but
// has some limitations. The state dir enables:
//
//   - persisted runner ID across restarts
//   - persisted adoption token across restarts
//
// The state directory will be created if it does not exist.
func WithStateDir(v string) Option {
	return func(r *Runner, cfg *config) error {
		r.stateDir = v
		return nil
	}
}

// WithAcceptTimeout sets a maximum amount of time to wait for a job before returning
// that one was not accepted.
func WithAcceptTimeout(dur time.Duration) Option {
	return func(r *Runner, cfg *config) error {
		r.acceptTimeout = dur
		return nil
	}
}

// WithLabels sets the labels for this runner.
func WithLabels(v map[string]string) Option {
	return func(r *Runner, cfg *config) error {
		cfg.labels = v
		return nil
	}
}
