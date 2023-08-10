// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// TestOnDemandRunnerConfig returns a valid project for tests.
func TestOnDemandRunnerConfig(t testing.T, src *pb.OnDemandRunnerConfig) *pb.OnDemandRunnerConfig {
	t.Helper()

	if src == nil {
		src = &pb.OnDemandRunnerConfig{
			PluginType: "docker",
		}
	}

	if src.TargetRunner == nil {
		src.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{
				Any: &pb.Ref_RunnerAny{},
			},
		}
	}

	require.NoError(t, mergo.Merge(src, &pb.OnDemandRunnerConfig{
		PluginType: "docker",
		OciUrl:     "hashicorp/waypoint-odr:latest",
	}))

	return src
}

// Type wrapper around the proto type so that we can add some methods.
type OnDemandRunnerConfig struct{ *pb.OnDemandRunnerConfig }

// ValidateOnDemandRunnerConfig validates the project structure.
func ValidateOnDemandRunnerConfig(p *pb.OnDemandRunnerConfig) error {
	return validationext.Error(validation.ValidateStruct(p,
		ValidateOnDemandRunnerConfigRules(p)...,
	))
}

// ValidateOnDemandRunnerConfigRules
func ValidateOnDemandRunnerConfigRules(p *pb.OnDemandRunnerConfig) []*validation.FieldRules {
	return []*validation.FieldRules{
		validation.Field(&p.PluginType, validation.Required),
		validation.Field(&p.PluginConfig, isPluginHcl(p)),
		validation.Field(&p.Name, validation.By(validatePathToken)),
		validation.Field(&p.TargetRunner, validation.By(func(i interface{}) error {
			targetRunner, _ := i.(*pb.Ref_Runner)
			if targetRunner == nil {
				return nil
			}
			switch runner := targetRunner.Target.(type) {
			case *pb.Ref_Runner_Labels:
				if len(runner.Labels.Labels) == 0 {
					return errors.New("profile targets a runner with labels, but does not specify any labels")
				}
			}
			return nil
		})),
	}
}

// ValidateUpsertOnDemandRunnerConfigRequest
func ValidateUpsertOnDemandRunnerConfigRequest(v *pb.UpsertOnDemandRunnerConfigRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Config, validation.Required),
		validationext.StructField(&v.Config, func() []*validation.FieldRules {
			return ValidateOnDemandRunnerConfigRules(v.Config)
		}),
	))
}

// ValidateGetOnDemandRunnerConfigRequest
func ValidateGetOnDemandRunnerConfigRequest(v *pb.GetOnDemandRunnerConfigRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Config, validation.Required),
	))
}

// ValidateDeleteOnDemandRunnerConfigRequest
func ValidateDeleteOnDemandRunnerConfigRequest(v *pb.DeleteOnDemandRunnerConfigRequest) error {
	return validationext.Error(validation.ValidateStruct(v,
		validation.Field(&v.Config, validation.Required),
	))
}

func isPluginHcl(p *pb.OnDemandRunnerConfig) validation.Rule {
	return validation.By(func(_ interface{}) error {
		if len(p.PluginConfig) == 0 {
			return nil
		}

		switch p.ConfigFormat {
		case pb.Hcl_HCL:
			_, diag := hclsyntax.ParseConfig(p.PluginConfig, "<waypoint-hcl>", hcl.Pos{})
			if diag.HasErrors() {
				return diag
			}

			return nil
		case pb.Hcl_JSON:
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
