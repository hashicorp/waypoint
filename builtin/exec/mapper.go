// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package exec

import (
	"github.com/hashicorp/waypoint/builtin/docker"
)

// DockerImageMapper maps a docker.Image to our Input structure.
func DockerImageMapper(src *docker.Image) *Input {
	return &Input{
		Data: map[string]*Input_Value{
			"DockerImageFull": {
				Value: &Input_Value_Text{
					Text: src.Name(),
				},
			},
			"DockerImageName": {
				Value: &Input_Value_Text{
					Text: src.Image,
				},
			},
			"DockerImageTag": {
				Value: &Input_Value_Text{
					Text: src.Tag,
				},
			},
		},
	}
}
