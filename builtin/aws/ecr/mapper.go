package ecr

import (
	"github.com/hashicorp/waypoint/builtin/docker"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ECRImageMapper maps a ecr.Image to a docker.Image structure.
func ECRImageMapper(src *Image) *docker.Image {
	return &docker.Image{
		Image: src.Image,
		Tag:   src.Tag,

		Location: &docker.Image_Registry{
			Registry: &emptypb.Empty{},
		},
	}
}
