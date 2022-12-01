package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
)

func (s *Service) UpsertProjectFromTemplate(
	ctx context.Context,
	req *pb.UpsertProjectFromTemplateRequest,
) (*pb.UpsertProjectFromTemplateResponse, error) {
	log := hclog.FromContext(ctx)

	// TODO: validate request

	// TODO(izaak): The architecture I want here is for the server to set up the project,
	// the runner to JUST create the new git repo, and then the server to finish up by
	// modifying the project with the new git repo. I don't see a great way to do this though -
	// we can only send the whole task to the runner and let it handle orchestration from there.

	template, err := s.state(ctx).ProjectTemplateGet(ctx, req.Template)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to get project template %q", req.Template.Name)
	}

	project := template.ProjectSettings
	project.Name = req.ProjectName

	// NOTE(izaak): It isn't great that the project exists for a while with the
	// template as it's datasource, but I think we need _a_ datasource or the remote
	// job won't execute.
	if err := s.state(ctx).ProjectPut(ctx, project); err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to upsert initial project",
			"project_name",
			req.ProjectName,
		)
	}

	job, err := s.QueueJob(ctx, &pb.QueueJobRequest{
		Job: &pb.Job{
			Application: &pb.Ref_Application{
				Project: req.ProjectName,
			},
			Operation: &pb.Job_TemplateProject{
				TemplateProject: &pb.Job_TemplateProjectOp{
					Req: req,
				},
			},
			TargetRunner: &pb.Ref_Runner{
				Target: &pb.Ref_Runner_Any{},
			},
		},
	})
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to queue project template job")
	}

	return &pb.UpsertProjectFromTemplateResponse{
		JobId: job.JobId,
	}, nil
}

func (s *Service) UpsertProjectTemplate(
	ctx context.Context,
	req *pb.UpsertProjectRequest,
) (*pb.UpsertApplicationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "upsert project template unimplemented")
}

func (s *Service) GetProjectTemplate(
	ctx context.Context,
	req *pb.GetProjectTemplateRequest,
) (*pb.GetProjectTemplateResponse, error) {
	// TODO(izaak): validate

	resp, err := s.state(ctx).ProjectTemplateGet(ctx, req.ProjectTemplate)

	return &pb.GetProjectTemplateResponse{
		ProjectTemplate: resp,
	}, err

}

func (s *Service) ListProjectTemplates(
	ctx context.Context,
	req *pb.ListProjectTemplatesRequest,
) (*pb.ListProjectTemplatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "list project templates unimplemented")
}

func (s *Service) DeleteProjectTemplate(
	ctx context.Context,
	req *pb.DeleteProjectTemplateRequest,
) (*pb.DeleteProjectTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "delete project template unimplemented")
}
