package core

import (
	"context"
	"fmt"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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
		_, err = a.client.SetMetadata(ctx, &pb.MetadataSetRequest{
			Scope: &pb.MetadataSetRequest_Project{
				Project: a.project.Ref(),
			},
			Value: &pb.MetadataSetRequest_FileChangeSignal{
				FileChangeSignal: projMD.FileChangeSignal,
			},
		})
		if err != nil {
			return err
		}
	}

	if appMD != nil {
		_, err = a.client.SetMetadata(ctx, &pb.MetadataSetRequest{
			Scope: &pb.MetadataSetRequest_Application{
				Application: a.Ref(),
			},
			Value: &pb.MetadataSetRequest_FileChangeSignal{
				FileChangeSignal: appMD.FileChangeSignal,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
