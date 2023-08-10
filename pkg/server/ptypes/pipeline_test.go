// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

		{
			"multiple root steps",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "root",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},

					"root2": {
						Name: "root2",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},
				}
			},
			`exactly one root`,
		},

		{
			"exec image required",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "root",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{},
						},
					},
				}
			},
			`image: cannot be blank`,
		},

		{
			"cycle",
			func(v *pb.Pipeline) {
				v.Steps = map[string]*pb.Pipeline_Step{
					"root": {
						Name: "root",
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},

					"A": {
						Name:      "A",
						DependsOn: []string{"B"},
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},

					"B": {
						Name:      "B",
						DependsOn: []string{"A"},
						Kind: &pb.Pipeline_Step_Exec_{
							Exec: &pb.Pipeline_Step_Exec{
								Image: "hashicorp/waypoint",
							},
						},
					},
				}
			},
			`one or more cycles`,
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
