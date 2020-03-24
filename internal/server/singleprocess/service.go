package singleprocess

import (
	"github.com/boltdb/bolt"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

// service implements the gRPC service for the server.
type service struct {
	db *bolt.DB
}

var _ pb.DevflowServer = (*service)(nil)
