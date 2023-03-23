// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidateGetWorkspaceRequest(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.GetWorkspaceRequest)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"ref is not set",
			func(v *pb.GetWorkspaceRequest) {
				v.Workspace = nil
			},
			"workspace: cannot be blank",
		},

		{
			"ref set, blank workspace value",
			func(v *pb.GetWorkspaceRequest) {
				v.Workspace = &pb.Ref_Workspace{Workspace: ""}
			},
			"workspace: cannot be blank",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			value := TestGetWorkspaceRequest(t, nil)
			if f := tt.Modify; f != nil {
				f(value)
			}

			err := ValidateGetWorkspaceRequest(value)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
