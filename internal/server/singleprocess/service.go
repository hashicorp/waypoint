package singleprocess

import (
	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-memdb"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:gen/"

// service implements the gRPC service for the server.
type service struct {
	// db is the persisted on-disk database used for historical data
	db *bolt.DB

	// inmem is our in-memory database that stores ephemeral data in an
	// easier-to-query way. Some of this data may be periodically persisted
	// but most of this data is meant to be lost when the process restarts.
	inmem *memdb.MemDB
}

// New returns a devflow server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(db *bolt.DB) (pb.DevflowServer, error) {
	// Initialize our DB
	if err := dbInit(db); err != nil {
		return nil, err
	}

	// Initialize our in-memory database
	inmem, err := memdbInit()
	if err != nil {
		return nil, err
	}

	return &service{db: db, inmem: inmem}, nil
}

var _ pb.DevflowServer = (*service)(nil)
