// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"
)

// Destroy will destroy all the physical resources for this app in the current
// configured workspace. If this returns an error, it is possible that the
// destroy is in a partial state.
func (a *App) Destroy(ctx context.Context) error {
	destroyers := []struct {
		DestroyFunc          func(context.Context) error
		DestroyWorkspaceFunc func(context.Context) error
	}{
		{
			a.destroyAllReleases,
			a.destroyReleaseWorkspace,
		},
		{
			a.destroyAllDeploys,
			a.destroyDeployWorkspace,
		},
	}

	// First we need to destroy all operations.
	for _, d := range destroyers {
		if err := d.DestroyFunc(ctx); err != nil {
			return err
		}
	}

	// Next we call the destroy workspace hooks.
	for _, d := range destroyers {
		if err := d.DestroyWorkspaceFunc(ctx); err != nil {
			return err
		}
	}

	return nil
}
