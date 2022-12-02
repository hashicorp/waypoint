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

	// TODO(izaak): A better architecture here would be for the server to set up the project,
	// the runner to JUST create the new git repo, and then the server to finish up by
	// modifying the project with the new git repo. I don't see a great way to do this though -
	// we can only send the whole task to the runner and let it handle orchestration from there.

	template, err := s.state(ctx).ProjectTemplateGet(ctx, req.Template)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to get project template %q", req.Template.Name)
	}

	project := template.ProjectSettings
	project.Name = req.ProjectName

	// The job won't execute without a valid hcl file
	project.WaypointHcl = []byte(`project = "tmp"`)

	// Insert a dummy datasource to pacify the job system. It will error submitting the remote job
	// if we don't have some kind of datasource here.

	// TODO: We need this to use our existing job-queueing logic, which assumes
	// a project already exists.
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
			Workspace: &pb.Ref_Workspace{
				Workspace: "default", // TODO: rethink
			},

			// TODO: We need this to use our existing job-queueing logic, which
			// assumes a datasource.
			DataSource: &pb.Job_DataSource{
				Source: &pb.Job_DataSource_Git{
					Git: &pb.Job_Git{
						Url: "https://github.com/izaaklauer/noop.git",
					},
				},
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
	req *pb.UpsertProjectTemplateRequest,
) (*pb.UpsertProjectTemplateResponse, error) {
	// TODO(izaak): validate
	log := hclog.FromContext(ctx)

	err := s.state(ctx).ProjectTemplatePut(ctx, req.ProjectTemplate)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to put project template", "name", req.ProjectTemplate.Name)
	}
	return &pb.UpsertProjectTemplateResponse{}, nil
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
	// TODO(izaak): validate
	log := hclog.FromContext(ctx)

	list, _, err := s.state(ctx).ProjectTemplateList(ctx, nil)
	if err != nil {
		return nil, hcerr.Externalize(log, err, "failed to list project templates")
	}
	return &pb.ListProjectTemplatesResponse{
		ProjectTemplates: list,
		Pagination:       nil,
	}, nil
}

func (s *Service) DeleteProjectTemplate(
	ctx context.Context,
	req *pb.DeleteProjectTemplateRequest,
) (*pb.DeleteProjectTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "delete project template unimplemented")
}
