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

// TestOndemandRunnerConfig returns a valid project for tests.
func TestOndemandRunnerConfig(t testing.T, src *pb.OndemandRunnerConfig) *pb.OndemandRunnerConfig {
	t.Helper()

	if src == nil {
		src = &pb.OndemandRunnerConfig{
			PluginType: "docker",
		}
	}

	require.NoError(t, mergo.Merge(src, &pb.OndemandRunnerConfig{
		PluginType: "docker",
		OciUrl:     "hashicorp/waypoint:stable",
	}))

	return src
}

// Type wrapper around the proto type so that we can add some methods.
type OndemandRunnerConfig struct{ *pb.OndemandRunnerConfig }

// ValidateOndemandRunnerConfig validates the project structure.
func ValidateOndemandRunnerConfig(p *pb.OndemandRunnerConfig) error {
	return validationext.Error(validation.ValidateStruct(p,
		ValidateOndemandRunnerConfigRules(p)...,
	))
}

// ValidateOndemandRunnerConfigRules
func ValidateOndemandRunnerConfigRules(p *pb.OndemandRunnerConfig) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&p.PluginType, validation.Required),
		validation.Field(&p.PluginConfig, isPluginHcl(p)),
	}
}

// ValidateUpsertOndemandRunnerConfigRequest
func ValidateUpsertOndemandRunnerConfigRequest(v *pb.UpsertOndemandRunnerConfigRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Config, validation.Required),
		validationext.StructField(&v.Config, func() []*validation.FieldRules {
			return ValidateOndemandRunnerConfigRules(v.Config)
		}),
	))
}

func isPluginHcl(p *pb.OndemandRunnerConfig) validation.Rule {
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
