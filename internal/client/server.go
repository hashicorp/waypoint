// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/hashicorp/waypoint/internal/server"
	"github.com/hashicorp/waypoint/pkg/protocolversion"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/singleprocess"
	"github.com/hashicorp/waypoint/pkg/serverclient"
)

// initServerClient will initialize a gRPC connection to the Waypoint server.
// This is called if a client wasn't explicitly given with WithClient.
//
// If a connection is successfully established, this will register connection
// closing and server cleanup with the Project cleanup function.
//
// This function will do one of two things:
//
//  1. If connection options were given, it'll attempt to connect to
//     an existing Waypoint server.
//
//  2. If WithLocal was specified and no connection addresses can be
//     found, this will spin up an in-memory server.
func (c *Project) initServerClient(ctx context.Context, cfg *config) (*grpc.ClientConn, error) {
	log := c.logger.Named("server")

	// If we're local, then connection is optional.
	opts := cfg.connectOpts
	if !c.noLocalServer {
		log.Trace("Local server may be created later - server credentials optional")
		opts = append(opts, serverclient.Optional())
	}

	// Connect. If we're local, this is set as optional so conn may be nil
	log.Info("attempting to source credentials and connect")
	conn, err := serverclient.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// If we established a connection
	if conn != nil {
		log.Debug("connection established with sourced credentials")
		c.cleanup(func() { conn.Close() })
		return conn, nil
	}

	// No connection, meaning we have to spin up a local server. This
	// can only be reached if we specified "Optional" to serverclient
	// which is only possible if we configured this client to support local
	// mode.
	log.Info("no server credentials found, using in-memory local server")
	return c.initLocalServer(ctx)
}

// initLocalServer starts the local server and configures p.client to
// point to it. This also configures p.localClosers so that all the
// resources are properly cleaned up on Close.
//
// If this returns an error, all resources associated with this operation
// will be closed, but the project can retry.
func (c *Project) initLocalServer(ctx context.Context) (*grpc.ClientConn, error) {
	log := c.logger.Named("server")
	c.localServer = true

	// We use this pointer to accumulate things we need to clean up
	// in the case of an error. On success we nil this variable which
	// doesn't close anything.
	var closers []io.Closer
	defer func() {
		for _, c := range closers {
			c.Close()
		}
	}()

	// TODO(mitchellh): path to this
	path := filepath.Join("data.db")
	log.Debug("opening local mode DB", "path", path)

	// Open our database
	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed opening boltdb path %s. Is another server already running against this db?: %w", path, err)
	}
	closers = append(closers, db)

	// Create our server
	impl, err := singleprocess.New(
		singleprocess.WithDB(db),
		singleprocess.WithLogger(log.Named("singleprocess")),
		singleprocess.WithSuperuser(), // local mode has no auth
	)
	if err != nil {
		return nil, err
	}

	// We listen on a random locally bound port
	// TODO(mitchellh): we should use Unix domain sockets if supported
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, err
	}
	closers = append(closers, ln)

	// Create a new cancellation context so we can cancel in the case of an error
	ctx, cancel := context.WithCancel(ctx)

	// Run the server
	log.Info("starting built-in server for local operations", "addr", ln.Addr().String())
	go server.Run(server.WithContext(ctx),
		server.WithLogger(log),
		server.WithGRPC(ln),
		server.WithImpl(impl),
	)

	// Connect to the local server
	conn, err := grpc.DialContext(ctx, ln.Addr().String(),
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(protocolversion.UnaryClientInterceptor(protocolversion.Current())),
		grpc.WithStreamInterceptor(protocolversion.StreamClientInterceptor(protocolversion.Current())),
	)
	if err != nil {
		cancel()
		return nil, err
	}

	// Setup our server config. The configuration is specifically set so
	// so that there is no advertise address which will disable the CEB
	// completely.
	client := pb.NewWaypointClient(conn)
	_, err = client.SetServerConfig(ctx, &pb.SetServerConfigRequest{
		Config: &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{
				{
					Addr: "",
				},
			},
		},
	})
	if err != nil {
		cancel()
		return nil, err
	}

	// Success, persist the closers
	cleanupClosers := closers
	closers = nil
	c.cleanup(func() {
		for _, c := range cleanupClosers {
			c.Close()
		}
	})
	_ = cancel // pacify vet lostcancel

	return conn, nil
}

// negotiateApiVersion negotiates the API version to use and validates
// that we are compatible to talk to the server.
func (c *Project) negotiateApiVersion(ctx context.Context) error {
	log := c.logger

	log.Trace("requesting version info from server")
	resp, err := c.client.GetVersionInfo(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	log.Info("server version info",
		"version", resp.Info.Version,
		"api_min", resp.Info.Api.Minimum,
		"api_current", resp.Info.Api.Current,
		"entrypoint_min", resp.Info.Entrypoint.Minimum,
		"entrypoint_current", resp.Info.Entrypoint.Current,
	)

	// Store the server version info
	c.serverVersion = resp.Info

	vsn, err := protocolversion.Negotiate(protocolversion.Current().Api, resp.Info.Api)
	if err != nil {
		return err
	}

	log.Info("negotiated api version", "version", vsn)
	return nil
}
