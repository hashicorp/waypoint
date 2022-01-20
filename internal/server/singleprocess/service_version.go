package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/hashicorp/waypoint/internal/protocolversion"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *service) GetVersionInfo(
	ctx context.Context,
	req *empty.Empty,
) (*pb.GetVersionInfoResponse, error) {
	return &pb.GetVersionInfoResponse{
		Info: protocolversion.Current(),
	}, nil
}
