// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package singleprocess

import (
	"context"

	"github.com/hashicorp/go-hclog"
	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	// Currently logged in user by default
	user := s.UserFromContext(ctx)

	// If we have a request, get that user
	if req.User != nil {
		var err error
		user, err = s.state(ctx).UserGet(ctx, req.User)
		if err != nil {
			return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get user")
		}
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func (s *Service) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (*pb.UpdateUserResponse, error) {
	if err := serverptypes.ValidateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// Get the user so that we don't overwrite fields that shouldn't be.
	user, err := s.state(ctx).UserGet(ctx, &pb.Ref_User{
		Ref: &pb.Ref_User_Id{
			Id: &pb.Ref_UserId{Id: req.User.Id},
		},
	})
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to get user in update")
	}

	// Update our writable fields
	user.Username = req.User.Username
	user.Display = req.User.Display

	// Write it
	if err := s.state(ctx).UserPut(ctx, user); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to update user")
	}

	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

func (s *Service) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteUserRequest(req); err != nil {
		return nil, err
	}

	if err := s.state(ctx).UserDelete(ctx, req.User); err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to delete user")
	}

	return &empty.Empty{}, nil
}

func (s *Service) ListUsers(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListUsersResponse, error) {
	users, err := s.state(ctx).UserList(ctx)
	if err != nil {
		return nil, hcerr.Externalize(hclog.FromContext(ctx), err, "failed to list users")
	}

	return &pb.ListUsersResponse{Users: users}, nil
}
