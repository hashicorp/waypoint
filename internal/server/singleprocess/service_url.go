package singleprocess

import (
	"context"
	"crypto/tls"
	"time"

	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func (s *service) initURLGuestAccount(acceptURLTerms bool) error {
	// Connect without auth to our API client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTimeout(10*time.Second))
	if s.urlConfig.APIInsecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// If it isn't insecure, then we have to specify that we're using TLS
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
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
			ServerId:  s.id,
			AcceptTos: acceptURLTerms,
		},
	)
	if err != nil {
		return err
	}

	s.urlConfig.APIToken = accountResp.Token
	return nil
}
