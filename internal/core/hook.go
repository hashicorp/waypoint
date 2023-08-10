// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"
	"os/exec"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/config"
)

// execHook executes the given hook. This will return any errors. This ignores
// on_failure configurations so this must be processed external.
func (a *App) execHook(ctx context.Context, log hclog.Logger, h *config.Hook) error {
	log.Debug("executing hook", "command", h.Command)

	// Get our writers
	stdout, stderr, err := a.UI.OutputWriters()
	if err != nil {
		log.Warn("error getting UI stdout/stderr", "err", err)
		return err
	}

	// Build our command
	cmd := exec.CommandContext(ctx, h.Command[0], h.Command[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Start
	if err := cmd.Start(); err != nil {
		log.Warn("error starting command", "err", err)
		return err
	}

	// Wait
	if err := cmd.Wait(); err != nil {
		L := log

		exiterr, ok := err.(*exec.ExitError)
		if ok {
			L = L.With("code", exiterr.ExitCode())
		}

		L.Warn("error running command", "err", err)
		return err
	}

	return nil
}
