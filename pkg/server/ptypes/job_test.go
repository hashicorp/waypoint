// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidateJob(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.Job)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"id is set",
			func(j *pb.Job) { j.Id = "nope" },
			"id: must be empty",
		},

		{
			"workspace is set",
			func(j *pb.Job) { j.Workspace = nil },
			"workspace: cannot be blank",
		},

		{
			"git: path good",
			func(j *pb.Job) {
				j.DataSource = &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Git{
						Git: &pb.Job_Git{
							Url:  "example.com",
							Path: "foo",
						},
					},
				}
			},
			"",
		},

		{
			"git: path has a ..",
			func(j *pb.Job) {
				j.DataSource = &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Git{
						Git: &pb.Job_Git{
							Url:  "example.com",
							Path: "../foo",
						},
					},
				}
			},
			"path: must not contain",
		},

		{
			"git: path is absolute",
			func(j *pb.Job) {
				j.DataSource = &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Git{
						Git: &pb.Job_Git{
							Url:  "example.com",
							Path: "/foo/bar",
						},
					},
				}
			},
			"path: must be relative",
		},

		{
			"git: path starts with ./",
			func(j *pb.Job) {
				j.DataSource = &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Git{
						Git: &pb.Job_Git{
							Url:  "example.com",
							Path: "./foo/bar",
						},
					},
				}
			},
			"path: relative path shouldn't",
		},

		{
			"git: path has repeating /",
			func(j *pb.Job) {
				j.DataSource = &pb.Job_DataSource{
					Source: &pb.Job_DataSource_Git{
						Git: &pb.Job_Git{
							Url:  "example.com",
							Path: "foo//bar",
						},
					},
				}
			},
			"path: path should not contain repeated",
		}}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			job := TestJobNew(t, nil)
			if f := tt.Modify; f != nil {
				f(job)
			}

			err := ValidateJob(job)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
