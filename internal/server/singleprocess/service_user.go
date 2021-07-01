package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
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
