// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ecr

import (
	"github.com/hashicorp/waypoint/builtin/docker"
)

// ECRImageMapper maps a ecr.Image to a docker.Image structure.
func ECRImageMapper(src *Image) *docker.Image {
	return &docker.Image{
		Image: src.Image,
		Tag:   src.Tag,
		Location: &docker.Image_Registry{
			Registry: &docker.Image_RegistryLocation{},
		},
	}
}

func DockerToEcrImageMapper(src *docker.Image) *Image {
	return &Image{
		Image: src.Image,
		Tag:   src.Tag,
	}
}
