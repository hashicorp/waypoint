// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package ceb contains the core logic for the custom entrypoint binary ("ceb").
//
// The CEB does not work on Windows.
package ceb

import (
	"context"
	"fmt"
	"github.com/hashicorp/waypoint/version"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/env"
	"github.com/hashicorp/waypoint/internal/pkg/gatedwriter"
	"github.com/hashicorp/waypoint/internal/plugin"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

const (
	envDeploymentId        = "WAYPOINT_DEPLOYMENT_ID"
	envServerAddr          = "WAYPOINT_SERVER_ADDR"
	envServerTls           = "WAYPOINT_SERVER_TLS"
	envServerTlsSkipVerify = "WAYPOINT_SERVER_TLS_SKIP_VERIFY"
	envCEBDisable          = "WAYPOINT_CEB_DISABLE"
	envCEBDisableExec      = "WAYPOINT_CEB_DISABLE_EXEC"
	envCEBServerRequired   = "WAYPOINT_CEB_SERVER_REQUIRED"
	envCEBToken            = "WAYPOINT_CEB_INVITE_TOKEN"

	// envLogLevel is the env var to set with the log level. This
	// env var matches the Waypoint CLI on purpose. This can be set on
	// the entrypoint process OR via app config (`waypoint config`).
	envLogLevel = "WAYPOINT_LOG_LEVEL"
)

const (
	DefaultPort = 5000
)

// CEB represents the state of a running CEB.
type CEB struct {
	id           string
	deploymentId string
	context      context.Context
	execIdx      int64
	execDisable  bool

	// stateCond and its associated locker are used to protect all the
	// state-prefixed fields. These state fields can be watched using this
	// cond for state changes in the CEB. Anyone waiting on stateCond should
	// also verify the context didn't cancel. The stateCond will be broadcasted
	// when the root context cancels.
	stateCond       *sync.Cond
	stateConfig     bool // config stream is connected
	stateChildReady bool // ready to start child command
	stateExit       bool // true when exiting

	// logger is the logger that should be used internally. Log messages
	// sent here with the proper log level (Info or higher) will also be
	// streamed to the server.
	logger hclog.Logger

	// logCh can be sent entries that will be sent to the server. If the
	// server connection is severed or too many entries are sent, some may
	// be dropped but the channel should always be consumed.
	logCh          chan *pb.LogBatch_Entry
	logGatedWriter *gatedwriter.Writer

	// clientMu must be held anytime reading/writing client. internally
	// you probably want to use waitClient() instead of this directly.
	clientMu   sync.Mutex
	clientCond *sync.Cond
	client     pb.WaypointClient

	// childSigCh can be sent signals which will be sent to the child command via kill(2).
	childSigCh chan os.Signal

	// childDoneCh is sent a value (incl. nil) when the child process exits.
	// This is not sent anything for restarts.
	childDoneCh <-chan error

	// childCmdCh can be sent new commands to restart the child process. New
	// commands will stop the old command first. Values sent here are coalesced
	// in case many changes are sent in a row.
	childCmdCh chan<- *exec.Cmd

	// childCmdBase is the base command to use for making any changes to the
	// child; use the copyCmd() function to copy this safetly to make changes.
	// Do not write to this directly.
	childCmdBase *exec.Cmd

	closedVal   *uint32
	cleanupFunc func()

	urlAgentMu     sync.Mutex
	urlAgentCtx    context.Context
	urlAgentCancel func()

	//---------------------------------------------------------------
	// Config sourcing

	// configPlugins is the mapping of config source type to launched plugin.
	configPlugins map[string]*plugin.Instance
}

// Run runs a CEB with the given options.
//
// This will run until the context is cancelled. If the context is cancelled,
// we will attempt to gracefully exit the underlying program and attempt to
// clean up all resources.
func Run(ctx context.Context, os ...Option) error {
	// Create our ID
	id, err := server.Id()
	if err != nil {
		return status.Errorf(codes.Internal,
			"failed to generate unique ID: %s", err)
	}

	// Defaults, initialization
	ceb := &CEB{
		id:            id,
		context:       ctx,
		configPlugins: map[string]*plugin.Instance{},
		stateCond:     sync.NewCond(&sync.Mutex{}),

		// for our atomic ops, we just use new() rather than addr operators (&)
		// so that we can be sure that the 64-bit alignment requirement is correct
		closedVal: new(uint32),
	}
	ceb.clientCond = sync.NewCond(&ceb.clientMu)
	defer ceb.Close()

	// Setup our default config sourcers.
	ceb.configPlugins = plugin.ConfigSourcers

	// Set our options
	var cfg config
	for _, o := range os {
		err := o(ceb, &cfg)
		if err != nil {
			return err
		}
	}

	// Setup our system logger
	ceb.initSystemLogger()

	// We replace the default hclog logger with our own so that all those
	// logs also appear in the log streaming output. We don't expect anything
	// to be using hclog.L() but this is there just in case it does.
	hclog.SetDefault(ceb.logger)

	// We're disabled also if we have no client set and the server address is empty.
	// This means we have nothing to connect to.
	cfg.disable = cfg.disable || (ceb.client == nil && cfg.ServerAddr == "")

	ceb.logger.Info("entrypoint starting",
		"deployment_id", ceb.deploymentId,
		"instance_id", ceb.id,
		"args", cfg.ExecArgs,
	)

	vsn := version.GetVersion()
	ceb.logger.Info("entrypoint version",
		"full_string", vsn.FullVersionNumber(true),
		"version", vsn.Version,
		"prerelease", vsn.VersionPrerelease,
		"metadata", vsn.VersionMetadata,
		"revision", vsn.Revision,
	)

	// Initialize our base child command. We do this before any server
	// connections and so on because if this fails we just want to fail fast
	// before any network activity.
	if err := ceb.initChildCmd(ctx, &cfg); err != nil {
		return err
	}

	// If we are enabled, initialize the CEB feature set.
	if err := ceb.init(ctx, &cfg, false); err != nil {
		return err
	}

	// Run our subprocess
	select {
	case err := <-ceb.childDoneCh:
		return err

	case <-ctx.Done():
		ceb.logger.Info("received cancellation request, waiting for child to exit")

		// Perform a state condition broadcast. Everyone blocking on state
		// changes should also be watching for stateExit or context cancellation.
		ceb.setState(&ceb.stateExit, true)

		// Wait for the child to end
		<-ceb.childDoneCh
	}

	return nil
}

// waitState waits for the given state boolean to go true. This boolean
// must be a pointer to a state field on ceb. This will also return if
// stateExit flips true. The return value notes whether we should exit.
func (ceb *CEB) waitState(state *bool, v bool) (exit bool) {
	ceb.stateCond.L.Lock()
	defer ceb.stateCond.L.Unlock()
	for *state != v && !ceb.stateExit {
		ceb.stateCond.Wait()
	}

	return ceb.stateExit
}

// setState sets the value of a state var on the ceb struct and broadcasts
// the condition variable.
func (ceb *CEB) setState(state *bool, v bool) {
	ceb.stateCond.L.Lock()
	defer ceb.stateCond.L.Unlock()
	*state = v
	ceb.stateCond.Broadcast()
}

// Close cleans up any resources created by the CEB and should be called
// to gracefully exit.
func (ceb *CEB) Close() error {
	// Only close ones
	if !atomic.CompareAndSwapUint32(ceb.closedVal, 0, 1) {
		return nil
	}

	if f := ceb.cleanupFunc; f != nil {
		f()
	}

	return nil
}

// closed returns true if Close was called
func (ceb *CEB) closed() bool {
	return atomic.LoadUint32(ceb.closedVal) != 0
}

// cleanup stacks cleanup functions to call when Close is called.
func (ceb *CEB) cleanup(f func()) {
	oldF := ceb.cleanupFunc
	ceb.cleanupFunc = func() {
		defer f()
		if oldF != nil {
			oldF()
		}
	}
}

// DeploymentId returns the deployment ID that this CEB represents.
func (ceb *CEB) DeploymentId() string {
	return ceb.deploymentId
}

type config struct {
	disable             bool
	cebPtr              *CEB
	ExecArgs            []string
	ServerAddr          string
	ServerRequired      bool
	ServerTls           bool
	ServerTlsSkipVerify bool
	InviteToken         string
	FileRewriteSignal   string

	URLServicePort int
}

type Option func(*CEB, *config) error

// WithEnvDefaults sets the configuration based on well-known accepted
// environment variables. If this is NOT called, then the environment variable
// based confiugration will be ignored.
func WithEnvDefaults() Option {
	return func(ceb *CEB, cfg *config) error {
		var port int
		portStr := os.Getenv("PORT")
		if portStr == "" {
			port = DefaultPort
			os.Setenv("PORT", strconv.Itoa(DefaultPort))
		} else {
			i, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("Invalid value of PORT: %s", err)
			}

			port = i
		}

		cfg.URLServicePort = port
		cfg.ServerAddr = os.Getenv(envServerAddr)

		var err error
		cfg.ServerRequired, err = env.GetBool(envCEBServerRequired, false)
		if err != nil {
			return err
		}

		cfg.ServerTls, err = env.GetBool(envServerTls, false)
		if err != nil {
			return err
		}

		cfg.ServerTlsSkipVerify, err = env.GetBool(envServerTlsSkipVerify, false)
		if err != nil {
			return err
		}

		cfg.InviteToken = os.Getenv(envCEBToken)

		cfg.disable, err = env.GetBool(envCEBDisable, false)
		if err != nil {
			return err
		}

		ceb.deploymentId = os.Getenv(envDeploymentId)

		ceb.execDisable, err = env.GetBool(envCEBDisableExec, false)
		if err != nil {
			return err
		}

		return nil
	}
}

// WithExec sets the binary and arguments for the child process that the
// ceb execs. If the first value is not absolute then we'll look for it on
// the PATH.
func WithExec(args []string) Option {
	return func(ceb *CEB, cfg *config) error {
		cfg.ExecArgs = args
		return nil
	}
}

// WithClient specifies the Waypoint client to use directly. This will
// override any env vars or any other form of client connection configuration.
func WithClient(client pb.WaypointClient) Option {
	return func(ceb *CEB, cfg *config) error {
		ceb.client = client
		return nil
	}
}

// withCEBValue is used by tests to get the CEB struct pointer from Run.
// This is a nasty pattern but its encapsulated behind test helpers.
func withCEBValue(cebCh chan<- *CEB) Option {
	return func(ceb *CEB, cfg *config) error {
		cebCh <- ceb
		return nil
	}
}
