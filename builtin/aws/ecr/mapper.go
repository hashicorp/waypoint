package ecr

import (
	"github.com/hashicorp/waypoint/builtin/docker"
)

// ECRImageMapper maps a ecr.Image to a docker.Image structure.
func ECRImageMapper(src *Image) *docker.Image {
	return &docker.Image{
		Image: src.Image,
		Tag:   src.Tag,
	}
}
