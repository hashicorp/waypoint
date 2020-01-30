package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/mitchellh/devflow/internal/builtin"
	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/mapper"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log := hclog.New(&hclog.LoggerOptions{
		Name:   "devflow",
		Level:  hclog.Trace,
		Color:  hclog.AutoColor,
		Output: os.Stderr,
	})

	log.Debug("decoding configuration")
	var config config.Config
	if err := hclsimple.DecodeFile("devflow.hcl", nil, &config); err != nil {
		log.Error("error decoding configuration", "error", err)
		return 1
	}

	if len(config.Apps) > 1 {
		log.Error("only one app is supported at this time")
		return 1
	}

	app := config.Apps[0]

	// Context for our CLI
	ctx := context.Background()

	// Load all our components
	source := &component.Source{App: app.Name, Path: "."}

	// Builder
	log.Info("loading component", "type", "builder", "name", app.Build.Type)
	raw, err := loadComponent(builtin.Builders, app.Build, ctx, log)
	if err != nil {
		log.Error("error loading builder", "error", err)
		return 1
	}
	builder := raw.(component.Builder)

	// Registry
	log.Info("loading component", "type", "registry", "name", app.Registry.Type)
	raw, err = loadComponent(builtin.Registries, app.Registry, ctx, log)
	if err != nil {
		log.Error("error loading registry", "error", err)
		return 1
	}
	registry := raw.(component.Registry)

	// Build
	buildFunc, err := mapper.NewFunc(builder.BuildFunc())
	if err != nil {
		log.Error("error preparing builder", "error", err)
		return 1
	}

	chain, err := buildFunc.Chain(builtin.Mappers, ctx, source, log)
	if err != nil {
		log.Error("error preparing builder", "error", err)
		return 1
	}
	log.Debug("function chain", "chain", chain.String())

	buildArtifact, err := chain.Call()
	if err != nil {
		log.Error("error running builder", "error", err)
		return 1
	}

	// Registry
	pushFunc, err := mapper.NewFunc(registry.PushFunc())
	if err != nil {
		log.Error("error preparing registry push", "error", err)
		return 1
	}

	chain, err = pushFunc.Chain(builtin.Mappers, ctx, source, log, buildArtifact)
	if err != nil {
		log.Error("error preparing registry push", "error", err)
		return 1
	}
	log.Debug("function chain", "chain", chain.String())

	artifact, err := chain.Call()
	if err != nil {
		log.Error("error pushing artifact to registry", "error", err)
		return 1
	}

	fmt.Printf("DONE: %#v\n", artifact)
	return 0
}

func loadComponent(
	f *mapper.Factory,
	c *config.Component,
	args ...interface{},
) (interface{}, error) {
	fn := f.Func(c.Type)
	if fn == nil {
		return nil, fmt.Errorf("unknown type: %q", c.Type)
	}

	raw, err := fn.Call(args...)
	if err != nil {
		return nil, err
	}

	diag := component.Configure(raw, c.Body, nil)
	if diag.HasErrors() {
		return nil, diag
	}

	return raw, nil
}
