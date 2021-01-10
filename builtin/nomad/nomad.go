package nomad

import (
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk"
	"github.com/hashicorp/waypoint/builtin/docker"
	"github.com/hashicorp/waypoint/builtin/files"
)

//go:generate protoc -I ../../.. --go_opt=plugins=grpc --go_out=../../.. waypoint/builtin/nomad/plugin.proto

// Options are the SDK options to use for instantiation for
// the Nomad plugin.
var Options = []sdk.Option{
	sdk.WithComponents(&Platform{}),
	sdk.WithMappers(DockerMapper, FilesMapper),
}

func FilesMapper(src *files.Files) *NomadSpec {
	parts := strings.Split(src.Path, "/")

	return &NomadSpec{
		Driver:   "exec",
		Command:  parts[len(parts)-1],
		Artifact: src.Path,
	}
}

func DockerMapper(src *docker.Image) *NomadSpec {
	return &NomadSpec{
		Driver: "docker",
		Image:  src.Image,
		Tag:    src.Tag,
	}
}
