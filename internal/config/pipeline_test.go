package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	//pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestPipeline(t *testing.T) {
	// Define various pipelines parsing use-cases
	cases := []struct {
		File     string
		Pipeline string
		Func     func(*testing.T, *Pipeline)
	}{
		{
			"pipeline.hcl",
			"dontexist",
			func(t *testing.T, c *Pipeline) {
				require.Nil(t, c)
			},
		},

		{
			"pipeline.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Id)
			},
		},

		{
			"pipeline_step.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Id)

				s, err := c.Step(nil)
				require.NoError(err)

				op := s.Operation()
				require.NotNil(t, op)

				var p testStepPluginConfig
				diag := op.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
				require.Equal("example.com/test", s.ImageURL)
			},
		},

		/*
			// TODO(briancain): make this test pass, add more tests to validate each step config
			{
				"pipeline_multi_step.hcl",
				"foo",
				func(t *testing.T, c *Pipeline) {
					require := require.New(t)

					require.NotNil(t, c)
					require.Equal("foo", c.Id)

					s, err := c.Step(nil)
					require.NoError(err)

					op := s.Operation()
					require.NotNil(t, op)

					var p testStepPluginConfig
					diag := op.Configure(&p, nil)
					if diag.HasErrors() {
						t.Fatal(diag.Error())
					}

					require.NotEmpty(t, p.config.Foo)
					require.Equal("example.com/test", s.ImageURL)
				},
			},
		*/
	}

	// Test all the cases
	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "pipelines", tt.File), &LoadOptions{
				Workspace: "default",
			})
			require.NoError(err)

			pipeline, err := cfg.Pipeline(tt.Pipeline, nil)
			require.NoError(err)

			tt.Func(t, pipeline)
		})
	}
}

// testStepPluginConfig implements component.Configurable to test that we
// decode HCL properly.
type testStepPluginConfig struct {
	config struct {
		Foo string `hcl:"foo,attr"`
		Bar string `hcl:"bar,optional"`
	}
}

func (p *testStepPluginConfig) Config() (interface{}, error) {
	return &p.config, nil
}
