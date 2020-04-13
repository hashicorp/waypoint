// Package ceb contains the core logic for the custom entrypoint binary ("ceb").
package ceb

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

const (
	envServerAddr = "DEVFLOW_SERVER_ADDR"
)

// CEB represents the state of a running CEB.
type CEB struct {
	logger hclog.Logger
	client pb.DevflowClient

	cleanupFunc func()
}

// New creates a new CEB with the given options.
//
// The context is only used to cancel any blocking initialization tasks
// for starting the CEB. Once this function returns, cancelling the given
// context has no effect.
func New(ctx context.Context, os ...Option) (*CEB, error) {
	// Defaults, initialization
	ceb := &CEB{
		logger: hclog.L(),
	}

	// Set our options
	var cfg config
	for _, o := range os {
		o(ceb, &cfg)
	}

	// Initialize our server connection
	if err := ceb.dialServer(ctx, &cfg); err != nil {
		return nil, status.Errorf(codes.Aborted,
			"failed to connect to server: %s", err)
	}

	return ceb, nil
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
	ServerAddr     string
	ServerInsecure bool
}

type Option func(*CEB, *config)

// WithEnvDefaults sets the configuration based on well-known accepted
// environment variables. If this is NOT called, then the environment variable
// based confiugration will be ignored.
func WithEnvDefaults() Option {
	return func(ceb *CEB, cfg *config) {
		cfg.ServerAddr = os.Getenv(envServerAddr)
	}
}
