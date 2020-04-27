package singleprocess

import (
	"github.com/boltdb/bolt"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/singleprocess/state"
)

//go:generate sh -c "protoc -I proto/ proto/*.proto --go_out=plugins=grpc:gen/"

// service implements the gRPC service for the server.
type service struct {
	// db is the persisted on-disk database used for historical data
	db *bolt.DB

	// state is the state management interface that provides functions for
	// safely mutating server state.
	state *state.State
}

// New returns a devflow server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(db *bolt.DB) (pb.DevflowServer, error) {
	// Initialize our DB
	if err := dbInit(db); err != nil {
		return nil, err
	}

	// Initialize our state
	st, err := state.New()
	if err != nil {
		return nil, err
	}

	return &service{db: db, state: st}, nil
}

var _ pb.DevflowServer = (*service)(nil)
