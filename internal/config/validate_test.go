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

			err = cfg.Validate()
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
			"pipeline_no_step.hcl",
			"'step' stanza",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "pipelines", tt.File), nil)
			require.NoError(err)

			err = cfg.Validate()
			if tt.Err == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Err)
		})
	}
}
