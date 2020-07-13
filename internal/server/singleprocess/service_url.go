package singleprocess

import (
	"context"
	"time"

	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"google.golang.org/grpc"
)

func (s *service) initURLGuestAccount() error {
	// Connect without auth to our API client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if s.urlConfig.APIInsecure {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(s.urlConfig.APIAddress, opts...)
	if err != nil {
		return err
	}
	client := wphznpb.NewWaypointHznClient(conn)

	// Request a guest account
	accountResp, err := client.RegisterGuestAccount(
		context.Background(),
		&wphznpb.RegisterGuestAccountRequest{
			ServerId: s.id,
		},
	)
	if err != nil {
		return err
	}

	s.urlConfig.APIToken = accountResp.Token
	return nil
}
