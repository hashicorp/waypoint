package exec

import (
	"github.com/hashicorp/waypoint/builtin/docker"
)

// DockerImageMapper maps a docker.Image to our Input structure.
func DockerImageMapper(src *docker.Image) *Input {
	return &Input{
		Data: map[string]*Input_Value{
			"DockerImageFull": &Input_Value{
				Value: &Input_Value_Text{
					Text: src.Name(),
				},
			},
			"DockerImageName": &Input_Value{
				Value: &Input_Value_Text{
					Text: src.Image,
				},
			},
			"DockerImageTag": &Input_Value{
				Value: &Input_Value_Text{
					Text: src.Tag,
				},
			},
		},
	}
}
