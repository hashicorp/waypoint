// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package serverinstall

import (
	"testing"

	"github.com/hashicorp/waypoint/internal/installutil"
)

func TestDeriveDefaultODRImage(t *testing.T) {
	tests := []struct {
		name        string
		serverImage string
		want        string
		wantErr     bool
	}{
		{
			"Short name (does not add docker.io/library)",
			"hashicorp/waypoint:latest",
			"hashicorp/waypoint-odr:latest",
			false,
		},
		{
			"Alpha",
			"ghcr.io/hashicorp/waypoint/alpha:latest",
			"ghcr.io/hashicorp/waypoint/alpha-odr:latest",
			false,
		},
		{
			"Custom registry with port (doesn't get confused by multiple colons)",
			"my.registry:5000/hashicorp/waypoint:latest",
			"my.registry:5000/hashicorp/waypoint-odr:latest",
			false,
		},
		{
			"Custom registry with port and no tag returns error (doesn't see the port as a tag)",
			"my.registry:5000/hashicorp/waypoint",
			"",
			true,
		},
		{
			"No tag returns an error",
			"hashicorp/waypoint",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := installutil.DeriveDefaultODRImage(tt.serverImage)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeriveDefaultODRImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeriveDefaultODRImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
