// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package statetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func init() {
	tests["trigger"] = []testFunc{
		TestTrigger,
	}
}

func TestTrigger(t *testing.T, factory Factory, restartF RestartFactory) {
	ctx := context.Background()
	t.Run("Get returns not found error if not exist", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		_, err := s.TriggerGet(ctx, &pb.Ref_Trigger{
			Id: "foo",
		})
		require.Error(err)
		require.Equal(codes.NotFound, status.Code(err))
	})

	t.Run("Put and Get", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		name := "pew"

		// Set
		err := s.TriggerPut(ctx, &pb.Trigger{
			Project: &pb.Ref_Project{Project: "p_test"},
			Name:    name,
			Id:      "t_test",
		})
		require.NoError(err)
		// Get id

		// Get exact by id
		{
			resp, err := s.TriggerGet(ctx, &pb.Ref_Trigger{
				Id: "t_test",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Update
		err = s.TriggerPut(ctx, &pb.Trigger{
			Project:     &pb.Ref_Project{Project: "p_test"},
			Description: "test",
			Name:        name,
			Id:          "t_test",
		})
		require.NoError(err)

		// Get exact by id
		{
			resp, err := s.TriggerGet(ctx, &pb.Ref_Trigger{
				Id: "t_test",
			})
			require.NoError(err)
			require.NotNil(resp)
			require.Equal(resp.Description, "test")
		}

		// Set with no proj returns an error
		err = s.TriggerPut(ctx, &pb.Trigger{
			Name: name,
			Id:   "test_test",
		})
		require.Error(err)
	})

	t.Run("Deletion", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Set
		err := s.TriggerPut(ctx, &pb.Trigger{
			Project: &pb.Ref_Project{Project: "p_test"},
			Id:      "t_test",
		})
		require.NoError(err)
		// Get id

		// Get exact by id
		{
			resp, err := s.TriggerGet(ctx, &pb.Ref_Trigger{
				Id: "t_test",
			})
			require.NoError(err)
			require.NotNil(resp)
		}

		// Delete it
		err = s.TriggerDelete(ctx, &pb.Ref_Trigger{
			Id: "t_test",
		})
		require.NoError(err)

		// It's gone
		{
			_, err := s.TriggerGet(ctx, &pb.Ref_Trigger{
				Id: "t_test",
			})
			require.Error(err)
		}
	})

	t.Run("Listing", func(t *testing.T) {
		require := require.New(t)

		s := factory(t)
		defer s.Close()

		// Create more for listing
		err := s.TriggerPut(ctx, &pb.Trigger{
			Project: &pb.Ref_Project{Project: "p_test"},
			Tags:    []string{"first"},
			Name:    "firsttest",
			Id:      "t_test",
		})
		require.NoError(err)

		err = s.TriggerPut(ctx, &pb.Trigger{
			Project: &pb.Ref_Project{Project: "test_project"},
			Application: &pb.Ref_Application{
				Project:     "test_project",
				Application: "a_test",
			},
			Name: "another_test",
			Id:   "t_test_part2",
		})
		require.NoError(err)

		err = s.TriggerPut(ctx, &pb.Trigger{
			Project: &pb.Ref_Project{Project: "test_project"},
			Tags:    []string{"test", "another"},
			Name:    "more_test",
			Id:      "t_test_part3",
		})
		require.NoError(err)

		// List all
		{
			resp, err := s.TriggerList(ctx, nil, nil, nil, nil)
			require.NoError(err)
			require.Len(resp, 3)
		}

		// List some
		{
			resp1, err := s.TriggerList(ctx, &pb.Ref_Workspace{Workspace: "default"}, &pb.Ref_Project{Project: "p_test"}, nil, nil)
			require.NoError(err)
			require.Len(resp1, 1)

			resp2, err := s.TriggerList(ctx, &pb.Ref_Workspace{Workspace: "default"}, &pb.Ref_Project{Project: "test_project"}, nil, nil)
			require.NoError(err)
			require.Len(resp2, 2)
		}

		// List none
		{
			resp, err := s.TriggerList(ctx, &pb.Ref_Workspace{Workspace: "production"}, nil, nil, nil)
			require.NoError(err)
			require.Len(resp, 0)
		}

		// List by workspace
		{
			resp, err := s.TriggerList(ctx, &pb.Ref_Workspace{Workspace: "default"}, nil, nil, nil)
			require.NoError(err)
			require.Len(resp, 3)
		}

		// List by project
		{
			resp, err := s.TriggerList(ctx, &pb.Ref_Workspace{Workspace: "default"}, &pb.Ref_Project{Project: "test_project"}, nil, nil)
			require.NoError(err)
			require.Len(resp, 2)
		}

		// List by application
		{
			// No app ref but project set means all apps, so 2
			resp, err := s.TriggerList(ctx, &pb.Ref_Workspace{Workspace: "default"},
				&pb.Ref_Project{Project: "test_project"}, &pb.Ref_Application{Project: "test_project", Application: "a_test"}, nil)
			require.NoError(err)
			require.Len(resp, 2)
		}

		// List by tag
		{
			resp, err := s.TriggerList(ctx, nil, nil, nil, []string{"test"})
			require.NoError(err)
			require.Len(resp, 1)
		}

	})
}
