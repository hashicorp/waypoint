package client

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/waypoint/internal/protocolversion"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// NegotiateApiVersion negotiates the API version to use and validates
// that we are compatible to talk to the server.
func NegotiateApiVersion(ctx context.Context, client pb.WaypointClient, log hclog.Logger) (*pb.VersionInfo, error) {

	log.Trace("requesting version info from server")
	resp, err := client.GetVersionInfo(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	log.Info("server version info",
		"version", resp.Info.Version,
		"api_min", resp.Info.Api.Minimum,
		"api_current", resp.Info.Api.Current,
		"entrypoint_min", resp.Info.Entrypoint.Minimum,
		"entrypoint_current", resp.Info.Entrypoint.Current,
	)

	vsn, err := protocolversion.Negotiate(protocolversion.Current().Api, resp.Info.Api)
	if err != nil {
		return resp.Info, err
	}

	log.Info("negotiated api version", "version", vsn)
	return resp.Info, nil
}
