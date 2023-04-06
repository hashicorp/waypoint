// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"github.com/hashicorp/opaqueany"
	"github.com/imdario/mergo"
	"github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Type wrapper around the proto type so that we can add some methods.
type Component struct{ *pb.Component }

// Match returns true if the component matches the given ref.
func (c *Component) Match(ref *pb.Ref_Component) bool {
	if c == nil || ref == nil {
		return false
	}

	return c.Type == ref.Type && c.Name == ref.Name
}

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
		Status: testStatus(t),
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
		Status: testStatus(t),
	}))

	return src
}

func TestValidDeployment(t testing.T, src *pb.Deployment) *pb.Deployment {
	t.Helper()

	if src == nil {
		src = &pb.Deployment{}
	}

	deployment, _ := opaqueany.New(&empty.Empty{})

	require.NoError(t, mergo.Merge(src, &pb.Deployment{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
		Status:     testStatus(t),
		Deployment: deployment,
	}))

	return src
}

func TestValidRelease(t testing.T, src *pb.Release) *pb.Release {
	t.Helper()

	if src == nil {
		src = &pb.Release{}
	}

	release, _ := opaqueany.New(&empty.Empty{})

	require.NoError(t, mergo.Merge(src, &pb.Release{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
		Status:  testStatus(t),
		Release: release,
	}))

	return src
}

func TestValidStatusReport(t testing.T, src *pb.StatusReport) *pb.StatusReport {
	t.Helper()

	if src == nil {
		src = &pb.StatusReport{}
	}

	require.NoError(t, mergo.Merge(src, &pb.StatusReport{
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
		DeprecatedResourcesHealth: []*pb.StatusReport_Health{
			{
				HealthStatus:  "READY",
				HealthMessage: "ready for requests",
			},
		},
	}))

	return src
}

func TestValidTrigger(t testing.T, src *pb.Trigger) *pb.Trigger {
	t.Helper()

	if src == nil {
		src = &pb.Trigger{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Trigger{
		Project: &pb.Ref_Project{
			Project: "p_test",
		},
		Application: &pb.Ref_Application{
			Application: "a_test",
			Project:     "p_test",
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: "default",
		},
		Operation: &pb.Trigger_Up{},
	}))

	return src
}

func TestValidTask(t testing.T, src *pb.Task) *pb.Task {
	t.Helper()

	if src == nil {
		src = &pb.Task{}
	}

	require.NoError(t, mergo.Merge(src, &pb.Task{}))

	return src
}

func testStatus(t testing.T) *pb.Status {
	pt := timestamppb.Now()

	return &pb.Status{
		State:        pb.Status_SUCCESS,
		StartTime:    pt,
		CompleteTime: pt,
	}
}
