package ptypes

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestRunner(t testing.T, src *pb.Runner) *pb.Runner {
	t.Helper()

	if src == nil {
		src = &pb.Runner{}
	}

	id, err := server.Id()
	require.NoError(t, err)

	require.NoError(t, mergo.Merge(src, &pb.Runner{
		Id: id,
	}))

	return src
}
