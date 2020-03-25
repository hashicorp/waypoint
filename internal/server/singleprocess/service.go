package singleprocess

import (
	"github.com/boltdb/bolt"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:gen/"

// service implements the gRPC service for the server.
type service struct {
	db *bolt.DB
}

// New returns a devflow server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(db *bolt.DB) (pb.DevflowServer, error) {
	// Initialize our DB
	if err := dbInit(db); err != nil {
		return nil, err
	}

	return &service{db: db}, nil
}

var _ pb.DevflowServer = (*service)(nil)
