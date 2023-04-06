// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestValidatePaginationRequest(t *testing.T) {
	cases := []struct {
		Name   string
		Modify func(*pb.PaginationRequest)
		Error  string
	}{
		{
			"valid empty PaginationRequest",
			nil,
			"",
		},

		{
			"valid previous page request",
			func(v *pb.PaginationRequest) {
				v.NextPageToken = "nextPageToken"
			},
			"",
		},

		{
			"valid next page request",
			func(v *pb.PaginationRequest) {
				v.PreviousPageToken = "previousPageToken"
			},
			"",
		},

		{
			"invalid - has both NextPageToken and PreviousPageToken",
			func(v *pb.PaginationRequest) {
				v.NextPageToken = "nextPageToken"
				v.PreviousPageToken = "previousPageToken"
			},
			"Only one of NextPageToken or PreviousPageToken can be set.",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)

			value := TestPaginationRequest(t, nil)
			if f := tt.Modify; f != nil {
				f(value)
			}

			err := ValidatePaginationRequest(value)
			if tt.Error == "" {
				require.NoError(err)
				return
			}

			require.Error(err)
			require.Contains(err.Error(), tt.Error)
		})
	}
}
