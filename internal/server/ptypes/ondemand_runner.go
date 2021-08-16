package ptypes

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/pkg/validationext"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// TestOndemandRunner returns a valid project for tests.
func TestOndemandRunner(t testing.T, src *pb.OndemandRunner) *pb.OndemandRunner {
	t.Helper()

	if src == nil {
		src = &pb.OndemandRunner{
			Id: "od_test",
		}
	}

	require.NoError(t, mergo.Merge(src, &pb.OndemandRunner{
		PluginType: "docker",
		OciUrl:     "hashicorp/waypoint:stable",
	}))

	return src
}

// Type wrapper around the proto type so that we can add some methods.
type OndemandRunner struct{ *pb.OndemandRunner }

// ValidateOndemandRunner validates the project structure.
func ValidateOndemandRunner(p *pb.OndemandRunner) error {
	return validationext.Error(validation.ValidateStruct(p,
		ValidateOndemandRunnerRules(p)...,
	))
}

// ValidateOndemandRunnerRules
func ValidateOndemandRunnerRules(p *pb.OndemandRunner) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&p.PluginType, validation.Required),
		validation.Field(&p.PluginConfig, isPluginHcl(p)),
	}
}

// ValidateUpsertOndemandRunnerRequest
func ValidateUpsertOndemandRunnerRequest(v *pb.UpsertOndemandRunnerRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.OndemandRunner, validation.Required),
		validationext.StructField(&v.OndemandRunner, func() []*validation.FieldRules {
			return ValidateOndemandRunnerRules(v.OndemandRunner)
		}),
	))
}

func isPluginHcl(p *pb.OndemandRunner) validation.Rule {
	return validation.By(func(_ interface{}) error {
		if len(p.PluginConfig) == 0 {
			return nil
		}

		switch p.ConfigFormat {
		case pb.Project_HCL:
			_, diag := hclsyntax.ParseConfig(p.PluginConfig, "<waypoint-hcl>", hcl.Pos{})
			if diag.HasErrors() {
				return diag
			}

			return nil
		case pb.Project_JSON:
			_, diag := hcljson.Parse(p.PluginConfig, "<waypoint-hcl>")
			if diag.HasErrors() {
				return diag
			}

			return nil
		default:
			return errors.New("unknown format")
		}
	})
}
