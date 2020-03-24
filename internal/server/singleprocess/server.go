package singleprocess

import (
	"github.com/boltdb/bolt"

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:gen/"

// New returns a devflow server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(db *bolt.DB) (pb.DevflowServer, error) {
	// Initialize our DB
	if err := dbInit(db); err != nil {
		return nil, err
	}

	return &service{db: db}, nil
}
