// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package core

import (
	"context"
	"fmt"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// ConfigSync writes all the app configuration in the waypoint.hcl file to
// the server.
func (a *App) ConfigSync(ctx context.Context) error {
	a.logger.Debug("evaluating config vars for syncing")
	vars, err := a.config.ConfigVars()
	if err != nil {
		return err
	}

	// If we have no vars then we don't want to round trip to the server.
	if len(vars) == 0 {
		a.logger.Debug("no file-based config vars, not syncing config")
		return nil
	}

	a.logger.Info("syncing config variables", "len", len(vars))
	if a.logger.IsDebug() {
		for _, v := range vars {
			a.logger.Debug("variable",
				"name", v.Name,
				"type", fmt.Sprintf("%T", v.Value))
		}
	}

	_, err = a.client.SetConfig(ctx, &pb.ConfigSetRequest{
		Variables: vars,
	})
	if err != nil {
		return err
	}

	// Sync the project and application metadata

	projMD, appMD := a.config.ConfigMetadata()

	if projMD != nil {
		resp, err := a.client.GetProject(ctx, &pb.GetProjectRequest{
			Project: a.project.Ref(),
		})
		if err != nil {
			return err
		}

		resp.Project.FileChangeSignal = projMD.FileChangeSignal

		_, err = a.client.UpsertProject(ctx, &pb.UpsertProjectRequest{
			Project: resp.Project,
		})
		if err != nil {
			return err
		}
	}

	if appMD != nil {
		_, err := a.client.UpsertApplication(ctx, &pb.UpsertApplicationRequest{
			Project:          a.project.Ref(),
			Name:             a.Ref().Application,
			FileChangeSignal: appMD.FileChangeSignal,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
