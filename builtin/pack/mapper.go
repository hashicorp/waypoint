package pack

import (
	"github.com/hashicorp/waypoint/builtin/docker"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PackImageMapper maps a pack.DockerImage to our Image structure.
//
// NOTE(mitchellh): the pack builder can probably just reuse the image
// from here but at the time of writing I was still building all the
// mapper subsystems so I wanted to test it out.
func PackImageMapper(src *DockerImage) *docker.Image {
	img := &docker.Image{
		Image: src.Image,
		Tag:   src.Tag,
	}

	if src.Remote {
		img.Location = &docker.Image_Registry{Registry: &emptypb.Empty{}}
	}

	return img
}
