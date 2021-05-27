package config

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/stretchr/testify/require"
)

func TestVariables_validateHcl(t *testing.T) {
	cases := []struct {
		File string
		Err  string
	}{
		{
			"valid.hcl",
			"",
		},
		{
			"invalid_type.hcl",
			"Invalid type specification",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "variables", tt.File), &LoadOptions{
				Workspace: "default",
			})
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

func TestVariables_decode(t *testing.T) {
	// TODO krantzinator: this can probably move under just validate, and
	// then Decode calls the validate function and if it passes, saves those in
	// *variables
	cases := []struct {
		File string
		Err  string
	}{
		{
			"valid_blocks.hcl",
			"",
		},
		{
			"duplicate_def.hcl",
			"Duplicate variable",
		},
		{
			"invalid_name.hcl",
			"Invalid variable name",
		},
		{
			"invalid_def.hcl",
			"Invalid default value",
		},
	}

	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "variables", tt.File), &LoadOptions{
				Workspace: "default",
			})
			require.NoError(err)

			schema, _ := gohcl.ImpliedBodySchema(&hclConfig{})
			content, diag := cfg.Body.Content(schema)
			require.False(diag.HasErrors())

			vars := &Variables{}
			for _, block := range content.Blocks {
				switch block.Type {
				case "variable":
					diag = vars.decodeVariableBlock(block)
				}
			}

			if tt.Err == "" {
				require.False(diag.HasErrors())
				return
			}

			require.True(diag.HasErrors())
			require.Contains(diag.Error(), tt.Err)
		})
	}
}
