// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aci

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidLocationReturnsError(t *testing.T) {
	err := validateLocationAvailable("abc", []string{"123", "foo"})
	require.Error(t, err)
}

func TestValidLocationReturnsNil(t *testing.T) {
	err := validateLocationAvailable("123", []string{"123", "foo"})
	require.NoError(t, err)
}

func TestValidateConfig(t *testing.T) {
	tests := map[string]struct {
		input Config
		valid bool
	}{
		"Valid Memory": {
			Config{
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 512, // max 4GB
				},
			},
			true,
		},
		"Error when two volume types are set": {
			Config{
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 512, // max 4GB
				},
				Volumes: []Volume{
					{
						Name:           "test",
						Path:           "/dfdf",
						AzureFileShare: &AzureFileShareVolume{Name: "sfsdf"},
						GitRepoVolume:  &GitRepoVolume{Repository: "sdfsadf"},
					},
				},
			},
			false,
		},
		"Error when no volume types are set": {
			Config{
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 512, // max 4GB
				},
				Volumes: []Volume{
					{
						Name: "test",
						Path: "/dfdf",
					},
				},
			},
			false,
		},
		"Valid when one volume type set": {
			Config{
				Location: "europe-north1",
				Capacity: &Capacity{
					Memory: 512, // max 4GB
				},
				Volumes: []Volume{
					{
						Name:           "test",
						Path:           "/dfdf",
						AzureFileShare: &AzureFileShareVolume{Name: "sfsdf"},
					},
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

func TestParseDockerImageWithServerReturnsServer(t *testing.T) {
	s := parseDockerServer("jacksonnic.azureacr.io/image:latest")
	require.Equal(t, "jacksonnic.azureacr.io", s)
}

func TestParseDockerImageWithCanonicalReturnsServer(t *testing.T) {
	s := parseDockerServer("docker.io/nicholasjackson/image:latest")
	require.Equal(t, "docker.io", s)
}

func TestParseDockerImageOfficialNameReturnsDefaultServer(t *testing.T) {
	s := parseDockerServer("image:latest")
	require.Equal(t, "docker.io", s)
}

func TestParseDockerNoServerReturnsDefaultServer(t *testing.T) {
	s := parseDockerServer("nicholasjackson/image:latest")
	require.Equal(t, "docker.io", s)
}
