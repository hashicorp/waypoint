package config

import (
	"path/filepath"
	"testing"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		{
			"pipeline_nested.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Name)

				steps := c.Steps
				require.Len(steps, 2)
				s := steps[0]

				var p testStepPluginConfig
				diag := s.Configure(&p, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, p.config.Foo)
				require.Equal("example.com/test", s.ImageURL)

				// This should be an embedded pipeline
				s2 := steps[1]
				embedPipeline := s2.Pipeline

				require.Equal("bar", embedPipeline.Name)
				require.Len(embedPipeline.Steps, 1)
				ps := embedPipeline.Steps[0]
				require.Equal("boo", ps.Name)

				var pt testStepPluginConfig
				diag = ps.Configure(&pt, nil)
				if diag.HasErrors() {
					t.Fatal(diag.Error())
				}

				require.NotEmpty(t, pt.config.Foo)
				require.Equal("nested", pt.config.Foo)
				require.Equal("example.com/test", ps.ImageURL)
			},
		},

		{
			"pipeline_step_workspace.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Name)

				steps := c.Steps
				require.Len(steps, 3)
				for _, step := range steps {
					expected := ""
					if step.Name == "testworkspace" {
						expected = "testworkspace"
					}
					require.Equal(expected, step.Workspace)
				}
			},
		},

		{
			"pipeline_step_workspace_nested.hcl",
			"foo",
			func(t *testing.T, c *Pipeline) {
				require := require.New(t)

				require.NotNil(t, c)
				require.Equal("foo", c.Name)

				steps := c.Steps
				require.Len(steps, 5)

				// expect step 0 to not have a pipeline or workspace set
				require.Nil(steps[0].Pipeline)
				require.Empty(steps[0].Workspace)

				// we expect step 1 to have the testworkspace workspace set
				require.Equal("testworkspace", steps[1].Workspace)

				// step 2 should have a pipeline
				pipe1 := steps[2]
				require.NotNil(pipe1.Pipeline)
				require.Len(pipe1.Pipeline.Steps, 1)
				require.Empty(pipe1.Pipeline.Steps[0].Workspace)

				// step 3 should have a pipeline
				pipe2 := steps[3]
				require.NotNil(pipe2.Pipeline)
				require.Len(pipe2.Pipeline.Steps, 3)
				require.Empty(pipe2.Pipeline.Steps[0].Workspace)
				require.Equal("dontoverride", pipe2.Pipeline.Steps[1].Workspace)
				require.Empty(pipe2.Pipeline.Steps[2].Workspace)
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
			"pipeline_invalid_step.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				_, err := c.PipelineProtos()
				require.Error(err)
				require.Equal(codes.Internal, status.Code(err))
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

		{
			"pipeline_nested_pipes.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				pipelines, err := c.PipelineProtos()
				require.NoError(err)
				require.Len(pipelines, 2)

				require.Equal(pipelines[0].Name, "nested")
				require.Equal(pipelines[1].Name, "foo")

				// validate nested pipeline was set as a ref on parent pipeline
				parentPipe := pipelines[1]
				require.Len(parentPipe.Steps, 2)
				require.Equal(parentPipe.Steps["test"].Name, "test")
				require.Equal(parentPipe.Steps["pipe"].Name, "pipe")
				pRef, ok := parentPipe.Steps["pipe"].Kind.(*pb.Pipeline_Step_Pipeline_)
				require.Equal(ok, true)
				require.NotNil(pRef)

				pipeOwner, ok := pRef.Pipeline.Ref.Ref.(*pb.Ref_Pipeline_Owner)
				require.Equal(ok, true)
				require.Equal(pipeOwner.Owner.Project.Project, "foo")
				require.Equal(pipeOwner.Owner.PipelineName, "nested")

				// validate nested pipeline was created
				require.Len(pipelines[0].Steps, 1)

				nestedStep := pipelines[0].Steps["test_nested"]
				require.Equal(nestedStep.Name, "test_nested")
			},
		},

		{
			"pipeline_nested_refs.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				pipelines, err := c.PipelineProtos()
				require.NoError(err)
				require.Len(pipelines, 2)

				require.Equal(pipelines[0].Name, "foo")
				require.Equal(pipelines[1].Name, "pipe2")

				// validate nested pipeline was set as a ref on parent pipeline
				parentPipe := pipelines[0]
				require.Len(parentPipe.Steps, 2)
				require.Equal(parentPipe.Steps["test"].Name, "test")
				require.Equal(parentPipe.Steps["pipe"].Name, "pipe")
				pRef, ok := parentPipe.Steps["pipe"].Kind.(*pb.Pipeline_Step_Pipeline_)
				require.Equal(ok, true)
				require.NotNil(pRef)

				pipeOwner, ok := pRef.Pipeline.Ref.Ref.(*pb.Ref_Pipeline_Owner)
				require.Equal(ok, true)
				require.Equal(pipeOwner.Owner.Project.Project, "foo")
				require.Equal(pipeOwner.Owner.PipelineName, "pipe2")
			},
		},
		{
			"pipeline_step_workspace.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				pipelines, err := c.PipelineProtos()
				require.NoError(err)
				require.Len(pipelines, 1)

				require.Equal(pipelines[0].Name, "foo")

				testStep := pipelines[0].Steps["testworkspace"]
				require.NotNil(testStep.Workspace)
				require.Equal("testworkspace", testStep.Workspace.Workspace)
			},
		},

		{
			"pipeline_step_workspace_nested.hcl",
			func(t *testing.T, c *Config) {
				require := require.New(t)

				pipelines, err := c.PipelineProtos()
				require.NoError(err)
				require.Len(pipelines, 3)

				// grab our pipes and test individually
				var fooPipe *pb.Pipeline
				var nestedPipe *pb.Pipeline
				var nestedWSPipe *pb.Pipeline
				for _, p := range pipelines {
					if p.Name == "foo" {
						fooPipe = p
					}
					if p.Name == "nested" {
						nestedPipe = p
					}
					if p.Name == "nested_workspace" {
						nestedWSPipe = p
					}
				}

				// The fooPipe has 5 steps, 3 of them are nested pipelines. The
				// other 2 should have 1 with no workspace (inherits default),
				// the other set as "testworkspace"
				require.Len(fooPipe.Steps, 5)
				require.Len(nestedPipe.Steps, 1)
				require.Len(nestedWSPipe.Steps, 3)
				require.Empty(fooPipe.Steps["test"].Workspace)
				require.Empty(fooPipe.Steps["normal"].Workspace)
				require.Equal("testworkspace", fooPipe.Steps["testworkspace"].Workspace.Workspace)

				// The nested pipe has 1 step, with no workspace specified
				require.Len(nestedPipe.Steps, 1)
				require.Empty(nestedPipe.Steps["test_nested"].Workspace)

				// The nestedWSPipe has 5 steps, 1 of them has a workspace
				// specified. The parent step has a workspace specified however
				// during parsing that value is not cascaded down to the
				// sub-steps; that's handled during step/job creation
				require.Len(nestedWSPipe.Steps, 3)
				require.Empty(nestedWSPipe.Steps["test_nested"].Workspace)
				require.Equal("dontoverride", nestedWSPipe.Steps["test_nested_dontoverride"].Workspace.Workspace)
				require.Empty(nestedWSPipe.Steps["test_nested_override"].Workspace)
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
