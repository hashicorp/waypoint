package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidatePipeline(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.Pipeline)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"no owner",
			func(v *pb.Pipeline) { v.Owner = nil },
			"Owner: cannot be blank",
		},

		{
			"project is blank",
			func(v *pb.Pipeline) {
				v.Owner = &pb.Pipeline_Project{
					Project: &pb.Ref_Project{Project: ""},
				}
			},
			"project: cannot be blank",
		},

		{
			"no steps",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{}
			},
			"steps: cannot be blank",
		},

		{
			"step name not set",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "",
					},
				}
			},
			"name: cannot be blank",
		},

		{
			"step name doesn't match key",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "bar",
					},
				}
			},
			`key "root" doesn't match`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			v := TestPipeline(t, nil)
			if f := tt.Modify; f != nil {
				f(v)
			}

			err := ValidatePipeline(v)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
