// Package ceb contains the core logic for the custom entrypoint binary ("ceb").
package ceb

import (
	"context"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

const (
	envDeploymentId   = "DEVFLOW_DEPLOYMENT_ID"
	envServerAddr     = "DEVFLOW_SERVER_ADDR"
	envServerInsecure = "DEVFLOW_SERVER_INSECURE"
)

// CEB represents the state of a running CEB.
type CEB struct {
	id       string
	logger   hclog.Logger
	context  context.Context
	client   pb.DevflowClient
	childCmd *exec.Cmd
	config   *pb.EntrypointConfig
	execIdx  int64

	cleanupFunc func()
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
		id:      id,
		logger:  hclog.L(),
		context: ctx,
	}
	defer ceb.Close()

	// Set our options
	var cfg config
	for _, o := range os {
		o(ceb, &cfg)
	}

	ceb.logger.Info("entrypoint starting",
		"deployment_id", cfg.DeploymentId,
		"instance_id", ceb.id,
		"args", cfg.ExecArgs,
	)

	// Initialize our server connection
	if err := ceb.dialServer(ctx, &cfg); err != nil {
		return status.Errorf(codes.Aborted,
			"failed to connect to server: %s", err)
	}

	// Initialize our command
	if err := ceb.initChildCmd(ctx, &cfg); err != nil {
		return status.Errorf(codes.Aborted,
			"failed to connect to server: %s", err)
	}

	// Get our configuration and start the long-running stream for it.
	if err := ceb.initConfigStream(ctx, &cfg); err != nil {
		return err
	}

	// Initialize our log stream
	// NOTE(mitchellh): at some point we want this to be configurable
	// but for now we're just going for it.
	if err := ceb.initLogStream(ctx, &cfg); err != nil {
		return err
	}

	// Run our subprocess
	errCh := ceb.execChildCmd(ctx)
	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		ceb.logger.Info("received cancellation request, gracefully exiting")
		ceb.childCmd.Process.Kill()
		<-errCh
	}

	return nil
}

// Close cleans up any resources created by the CEB and should be called
// to gracefully exit.
func (ceb *CEB) Close() error {
	if f := ceb.cleanupFunc; f != nil {
		f()
	}

	return nil
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

type config struct {
	ExecArgs       []string
	DeploymentId   string
	ServerAddr     string
	ServerInsecure bool
}

type Option func(*CEB, *config)

// WithEnvDefaults sets the configuration based on well-known accepted
// environment variables. If this is NOT called, then the environment variable
// based confiugration will be ignored.
func WithEnvDefaults() Option {
	return func(ceb *CEB, cfg *config) {
		cfg.DeploymentId = os.Getenv(envDeploymentId)
		cfg.ServerAddr = os.Getenv(envServerAddr)
		cfg.ServerInsecure = os.Getenv(envServerInsecure) != ""
	}
}

// WithExec sets the binary and arguments for the child process that the
// ceb execs. If the first value is not absolute then we'll look for it on
// the PATH.
func WithExec(args []string) Option {
	return func(ceb *CEB, cfg *config) {
		cfg.ExecArgs = args
	}
}
