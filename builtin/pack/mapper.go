package pack

import (
	"github.com/hashicorp/waypoint/builtin/docker"
)

// PackImageMapper maps a pack.DockerImage to our Image structure.
//
// NOTE(mitchellh): the pack builder can probably just reuse the image
// from here but at the time of writing I was still building all the
// mapper subsystems so I wanted to test it out.
func PackImageMapper(src *DockerImage) *docker.Image {
	return &docker.Image{
		Image: src.Image,
		Tag:   src.Tag,
	}
}
