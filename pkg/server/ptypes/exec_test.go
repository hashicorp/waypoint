// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidateExecStreamRequestStart(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.ExecStreamRequest_Start)
		Error  string
	}{
		{
			"valid",
			nil,
			"",
		},

		{
			"args must not be blank",
			func(v *pb.ExecStreamRequest_Start) {
				v.Args = []string{}
			},
			"args: cannot be blank",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			value := TestExecStreamRequestStart(t, nil)
			if f := tt.Modify; f != nil {
				f(value)
			}

			err := ValidateExecStreamRequestStart(value)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
