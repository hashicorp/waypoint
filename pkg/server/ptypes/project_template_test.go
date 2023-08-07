package ptypes

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"testing"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidateCreateProjectTemplateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.CreateProjectTemplateRequest
		wantErr bool
	}{
		{
			name: "minimum valid request",
			req: &pb.CreateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{
					Name: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "enforces name",
			req: &pb.CreateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{},
			},
			wantErr: true,
		},
		{
			name: "Inherits base validator rules (example: name length)",
			req: &pb.CreateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{
					Name: string(make([]byte, PROJECT_TEMPLATE_NAME_LENGTH+1)),
				},
			},
			wantErr: true,
		},
		{
			name: "enforces name to not be empty",
			req: &pb.CreateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{
					Name: " "},
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateProjectTemplateRequest(tt.req)
			if err == nil && tt.wantErr {
				t.Errorf("expected error in ValidateCreateProjectTemplateRequest() but got none")
			}

			if err != nil && !tt.wantErr {
				t.Errorf("ValidateCreateProjectTemplateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestValidateUpdateProjectTemplateRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.UpdateProjectTemplateRequest
		wantErr bool
	}{
		{
			name: "Fails if no name or ID is given",
			req: &pb.UpdateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{},
			},
			wantErr: true,
		},
		{
			name: "OK with just name",
			req: &pb.UpdateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{
					Name: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "OK with just ID",
			req: &pb.UpdateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{
					Id: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "Also runs base project template validator",
			req: &pb.UpdateProjectTemplateRequest{
				ProjectTemplate: &pb.ProjectTemplate{
					Id: string(make([]byte, PROJECT_TEMPLATE_ID_LENGTH+1)),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateProjectTemplateRequest(tt.req)

			if err == nil && tt.wantErr {
				t.Errorf("expected error in ValidateUpdateProjectTemplateRequest() but got none")
			}

			if err != nil && !tt.wantErr {
				t.Errorf("ValidateUpdateProjectTemplateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateProjectTemplate(t *testing.T) {
	tests := []struct {
		name    string
		pt      *pb.ProjectTemplate
		wantErr bool
	}{
		{
			name:    "Fine with empty values",
			pt:      &pb.ProjectTemplate{},
			wantErr: false,
		},
		{
			name: "If name is set, enforces length limits",
			pt: &pb.ProjectTemplate{
				Name: string(make([]byte, PROJECT_TEMPLATE_NAME_LENGTH+1)),
			},
			wantErr: true,
		},
		{
			name: "Validates nexted TFC-related lengths",
			pt: &pb.ProjectTemplate{
				TerraformNocodeModule: &pb.ProjectTemplate_TerraformNocodeModule{
					Source:  "", // Empty string shouldn't be allowed
					Version: "0.0.1",
				},
			},
			wantErr: true,
		},
		{
			name: "Tag lengths",
			pt: &pb.ProjectTemplate{
				Tags: []string{string(make([]byte, PROJECT_TEMPLATE_TAG_LENGTH+1))},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateStruct(tt.pt, validateProjectTemplateRules(tt.pt)...)
			if err == nil && tt.wantErr {
				t.Errorf("expected error in validation but got none")
			}
			if err != nil && !tt.wantErr {
				t.Errorf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
