package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
				require.Equal("foo", c.Name)
			},
		},

		{
			"pipeline_step.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Name)

				steps := c.Steps
				s := steps[0]

				var p testStepPluginConfig
				diag := s.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
				require.Equal("example.com/test", s.ImageURL)
			},
		},

		{
			"pipeline_multi_step.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Name)

				steps := c.Steps
				require.Len(steps, 3)

				s := steps[0]

				var p testStepPluginConfig

				diag := s.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
				require.Equal("example.com/test", s.ImageURL)
				require.Equal("qubit", p.config.Foo)

				s2 := steps[1]

				diag = s2.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
				require.Equal("example.com/second", s2.ImageURL)
				require.Equal("few", p.config.Foo)
				require.Equal("bar", p.config.Bar)

				s3 := steps[2]

				diag = s3.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
				require.Equal("example.com/different", s3.ImageURL)
				require.Equal("food", p.config.Foo)
				require.Equal("drink", p.config.Bar)
				require.Len(s3.DependsOn, 1)
				require.Equal("zero", s3.DependsOn[0])
			},
		},
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

// TestPipelineProtos will test that given a config, we can translate a Pipeline
// HCL raw config into a Pipeline Proto that could be used to upsert the latest
// config into the Waypoint database.
func TestPipelineProtos(t *testing.T) {
	cases := []struct {
		File string
		Func func(*testing.T, *Config)
	}{
		// TODO verify a step exec was set properly and the fields are set on the proto
		{
			"pipeline_exec_step.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				pipelines, err := c.PipelineProtos()
				require.NoError(err)
				require.Len(pipelines, 1)

				require.Equal(pipelines[0].Name, "foo")
			},
		},

		{
			"pipeline_exec_step_many.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				pipelines, err := c.PipelineProtos()
				require.NoError(err)
				require.Len(pipelines, 2)

				require.Equal(pipelines[0].Name, "foo")
				require.Equal(pipelines[1].Name, "bar")
			},
		},
	}

	// Test all the cases
	for _, tt := range cases {
		t.Run(tt.File, func(t *testing.T) {
			require := require.New(t)

			cfg, err := Load(filepath.Join("testdata", "pipelines", tt.File), &LoadOptions{
				Workspace: "default",
			})
			require.NoError(err)

			tt.Func(t, cfg)
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
