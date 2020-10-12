package exec

import (
	"github.com/hashicorp/waypoint/builtin/docker"
)

// DockerImageMapper maps a docker.Image to our Input structure.
func DockerImageMapper(src *docker.Image) *Input {
	return &Input{}
}
