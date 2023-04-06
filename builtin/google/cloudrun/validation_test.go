// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudrun

import (
	"testing"

	"github.com/stretchr/testify/require"
	run "google.golang.org/api/run/v1"
)

func TestValidateImageReturnsErrorOnInvalidImageName(t *testing.T) {
	err := validateImageName("foo")
	require.Error(t, err)
}

func TestValidateImageReturnsErrorOnInvalidRegistry(t *testing.T) {
	err := validateImageName("foo/proj/image")
	require.Error(t, err)
}

func TestValidateImageReturnsErrorOnInvalidArtifactRegistry(t *testing.T) {
	err := validateImageName("FOO-docker.pkg.dev/waypoint-286812/foo/bar")
	require.Error(t, err)
}

func TestValidateImageReturnsNoErrorOnValidArtifactRegistry(t *testing.T) {
	err := validateImageName("europe-north1-docker.pkg.dev/waypoint-286812/foo/bar")
	require.NoError(t, err)
}

var locations = []*run.Location{
	{LocationId: "asia-east1"},
	{LocationId: "asia-northeast1"},
}

func TestValidateLocationAvailableReturnsErrorWhenLocationNotAvailable(t *testing.T) {
	err := validateLocationAvailable("badlocation", locations)
	require.Error(t, err)
}

func TestValidateLocationAvailableReturnsNoErrorWhenLocationAvailable(t *testing.T) {
	err := validateLocationAvailable("asia-east1", locations)
	require.NoError(t, err)
}

func TestConfigValidation(t *testing.T) {
	tests := map[string]struct {
		input Config
		valid bool
	}{
		"Valid Memory": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 128, // max 4GB
				},
			},
			true,
		},
		"Invalid Memory Value Too High": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 5000, // max 4GB
				},
			},
			false,
		},
		"Invalid Memory Value Too Low": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 50, // max 4GB
				},
			},
			false,
		},
		"Valid CPU Count": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					CPUCount: 2, // max 2
				},
			},
			true,
		},
		"CPU Count greater than max of 2": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					CPUCount: 3, // max 2
				},
			},
			false,
		},
		"Request Timeout valid": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					RequestTimeout: 300,
				},
			},
			true,
		},
		"Request Timeout greater than max": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					RequestTimeout: 901, // max 900
				},
			},
			false,
		},
		"Max requests per container valid": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					MaxRequestsPerContainer: 80,
				},
			},
			true,
		},
		"Max requests per container less than 0": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				Capacity: &Capacity{
					MaxRequestsPerContainer: -1,
				},
			},
			false,
		},
		"Autoscaling max valid": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				AutoScaling: &AutoScaling{
					Max: 0,
				},
			},
			true,
		},
		"Autoscaling max invalid": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				AutoScaling: &AutoScaling{
					Max: -1,
				},
			},
			false,
		},
		"Egress invalid": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				VPCAccess: &VPCAccess{
					Egress: "invalid value for egress",
				},
			},
			false,
		},
		"Egress valid": {
			Config{
				Project:  "waypoint-286812",
				Location: "europe-north1",
				VPCAccess: &VPCAccess{
					Egress: "all",
				},
			},
			true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateConfig(tc.input)

			if tc.valid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}

func TestValidateConfigReturnsPrettyError(t *testing.T) {
	c := Config{
		Project:  "waypoint-286812",
		Location: "europe-north1",
		Capacity: &Capacity{
			Memory:                  5000, // max 4GB
			CPUCount:                4,
			MaxRequestsPerContainer: -1,
			RequestTimeout:          1000,
		},
		AutoScaling: &AutoScaling{
			Max: -1,
		},
	}

	err := validateConfig(c)

	require.Error(t, err)
	require.Contains(t, err.Error(), ErrInvalidMemoryValue.Error())
	require.Contains(t, err.Error(), ErrInvalidCPUCount.Error())
	require.Contains(t, err.Error(), ErrInvalidRequestTimetout.Error())
	require.Contains(t, err.Error(), ErrInvalidMaxRequests.Error())
	require.Contains(t, err.Error(), ErrInvalidAutoscalingMax.Error())
}

func TestConfigSetReturnsErrorOnInvalidConfig(t *testing.T) {
	c := Config{
		Project:  "waypoint-286812",
		Location: "europe-north1",
		Capacity: &Capacity{
			Memory: 5000, // max 4096 (4GB)
		},
	}

	p := Platform{}

	err := p.ConfigSet(c)

	require.Error(t, err)
}

func TestConfigSetReturnsErrorOnInvalidInterface(t *testing.T) {
	p := Platform{}

	err := p.ConfigSet(struct{}{})

	require.Error(t, err)
}
