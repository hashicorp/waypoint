package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func (s *service) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
	// Currently logged in user by default
	user := s.userFromContext(ctx)

	// If we have a request, get that user
	if req.User != nil {
		var err error
		user, err = s.state.UserGet(req.User)
		if err != nil {
			return nil, err
		}
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

func (s *service) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (*pb.UpdateUserResponse, error) {
	if err := serverptypes.ValidateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// Get the user so that we don't overwrite fields that shouldn't be.
	user, err := s.state.UserGet(&pb.Ref_User{
		Ref: &pb.Ref_User_Id{
			Id: &pb.Ref_UserId{Id: req.User.Id},
		},
	})
	if err != nil {
		return nil, err
	}

	// Update our writable fields
	user.Username = req.User.Username
	user.Display = req.User.Display

	// Write it
	if err := s.state.UserPut(user); err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

func (s *service) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*empty.Empty, error) {
	if err := serverptypes.ValidateDeleteUserRequest(req); err != nil {
		return nil, err
	}

	if err := s.state.UserDelete(req.User); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *service) ListUsers(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListUsersResponse, error) {
	users, err := s.state.UserList()
	if err != nil {
		return nil, err
	}

	return &pb.ListUsersResponse{Users: users}, nil
}
