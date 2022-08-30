package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	cases := []struct {
		File string
		Err  string
	}{
		{
			"valid.hcl",
			"",
		},
		{
			"no_build.hcl",
			"'build' stanza",
		},

		// This isn't an error because we want to catch this at runtime.
		{
			"build_no_use.hcl",
			"",
		},

		{
			"build_scoped.hcl",
			"",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "validate", tt.File), nil)
			require.NoError(err)

			_, err = cfg.Validate()
			if tt.Err == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Err)
		})
	}
}

func TestConfigValidatePipelines(t *testing.T) {
	cases := []struct {
		File string
		Err  string
	}{
		{
			"pipeline_step.hcl",
			"",
		},
		{
			"pipeline_multi_step.hcl",
			"",
		},
		{
			"pipeline_nested_pipes.hcl",
			"",
		},
		{
			"pipeline_nested_refs.hcl",
			"",
		},
		{
			"pipeline_no_step.hcl",
			"'step' stanza",
		},
		{
			"pipeline_dupe_name.hcl",
			"'pipeline' stanza",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "pipelines", tt.File), nil)
			require.NoError(err)

			_, err = cfg.Validate()
			if tt.Err == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Err)

		})
	}
}

func TestConfigValidatePipelineSteps(t *testing.T) {
	cases := []struct {
		File string
		Err  string
	}{
		{
			"pipeline_step.hcl",
			"",
		},
		{
			"pipeline_multi_step.hcl",
			"",
		},
		{
			"pipeline_nested_pipes.hcl",
			"",
		},
		{
			"pipeline_nested_refs.hcl",
			"",
		},
		{
			"pipeline_step_no_use.hcl",
			"step stage with a default 'use'",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "pipelines", tt.File), nil)
			require.NoError(err)

			for _, name := range cfg.Pipelines() {
				pipeline, _ := cfg.Pipeline(name, nil)
				err = pipeline.Validate()
				if tt.Err == "" {
					require.NoError(err)
					return
				}
				require.Error(err)
				require.Contains(err.Error(), tt.Err)
			}
		})
	}
}
