package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/mitchellh/devflow/internal/component"
	"github.com/mitchellh/devflow/internal/config"
	"github.com/mitchellh/devflow/internal/core"
	"github.com/mitchellh/devflow/internal/datadir"
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

	// Context for our CLI
	ctx := context.Background()

	log.Debug("decoding configuration")
	var cfg config.Config
	if err := hclsimple.DecodeFile("devflow.hcl", nil, &cfg); err != nil {
		log.Error("error decoding configuration", "error", err)
		return 1
	}

	// Setup our directory
	log.Debug("preparing project directory", "path", ".devflow")
	projDir, err := datadir.NewProject(".devflow")
	if err != nil {
		log.Error("error preparing data directory", "error", err)
		return 1
	}

	// Create our project
	proj, err := core.NewProject(ctx,
		core.WithLogger(log),
		core.WithConfig(&cfg),
		core.WithDataDir(projDir),
	)
	if err != nil {
		log.Error("failed to create project", "error", err)
		return 1
	}

	// NOTE(mitchellh): temporary restriction
	if len(cfg.Apps) != 1 {
		log.Error("only one app is supported at this time")
		return 1
	}

	// Get our app
	app, err := proj.App(cfg.Apps[0].Name)
	if err != nil {
		log.Error("failed to initialize app", "error", err)
		return 1
	}

	// Build
	fmt.Fprintf(os.Stdout, "==> Building\n")
	buildArtifact, err := app.Build(ctx)
	if err != nil {
		log.Error("error running builder", "error", err)
		return 1
	}

	var pushedArtifact component.Artifact

	if app.Registry != nil {
		fmt.Fprintf(os.Stdout, "==> Pushing artifact\n")
		pushedArtifact, err = app.Push(ctx, buildArtifact)
		if err != nil {
			log.Error("error pushing artifact to registry", "error", err)
			return 1
		}
	} else {
		pushedArtifact = buildArtifact
	}

	fmt.Fprintf(os.Stdout, "==> Deploying\n")
	deployment, err := app.Deploy(ctx, pushedArtifact)
	if err != nil {
		log.Error("error deploying", "error", err)
		return 1
	}

	fmt.Fprintf(os.Stdout, "%s\n", deployment.String())
	return 0
}
