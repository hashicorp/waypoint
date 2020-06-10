package ptypes

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func TestValidBuild(t testing.T, src *pb.Build) *pb.Build {
	t.Helper()

	if src == nil {
		src = &pb.Build{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Build{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
	}))

	return src
}

func TestValidArtifact(t testing.T, src *pb.PushedArtifact) *pb.PushedArtifact {
	t.Helper()

	if src == nil {
		src = &pb.PushedArtifact{}
	}

	require.NoError(t, mergo.Merge(src, &pb.PushedArtifact{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
	}))

	return src
}

func TestValidDeployment(t testing.T, src *pb.Deployment) *pb.Deployment {
	t.Helper()

	if src == nil {
		src = &pb.Deployment{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Deployment{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
	}))

	return src
}
