package singleprocess

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/version"
)

//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//
// Protocol Versions
//
// These define the protocol versions supported by the server. You must be
// VERY THOUGHTFUL when modifying these values. Please read and re-read our
// upgrade policy to understand how these values work.
//
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
const (
	protocolVersionApiCurrent        uint32 = 1
	protocolVersionApiMin                   = 1
	protocolVersionEntrypointCurrent uint32 = 1
	protocolVersionEntrypointMin            = 1
)

func (s *service) GetVersionInfo(
	ctx context.Context,
	req *empty.Empty,
) (*pb.GetVersionInfoResponse, error) {
	return &pb.GetVersionInfoResponse{
		Info: &pb.VersionInfo{
			Api: &pb.VersionInfo_ProtocolVersion{
				Current: protocolVersionApiCurrent,
				Minimum: protocolVersionApiMin,
			},

			Entrypoint: &pb.VersionInfo_ProtocolVersion{
				Current: protocolVersionEntrypointCurrent,
				Minimum: protocolVersionEntrypointMin,
			},

			Version: version.GetVersion().VersionNumber(),
		},
	}, nil
}
