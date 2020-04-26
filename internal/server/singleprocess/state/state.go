// Package state manages the state that the singleprocess server has, providing
// operations to mutate that state safely as needed.
package state

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

// State is the primary API for state mutation for the server.
type State struct {
	// inmem is our in-memory database that stores ephemeral data in an
	// easier-to-query way. Some of this data may be periodically persisted
	// but most of this data is meant to be lost when the process restarts.
	inmem *memdb.MemDB
}

// New initializes a new State store.
func New() (*State, error) {
	// Create the in-memory DB.
	inmem, err := memdb.NewMemDB(stateStoreSchema())
	if err != nil {
		return nil, fmt.Errorf("Failed setting up state store: %s", err)
	}

	return &State{inmem: inmem}, nil
}

// Close should be called to gracefully close any resources.
func (s *State) Close() error {
	// Nothing for now, but we expect to do things one day.
	return nil
}

// schemaFn is an interface function used to create and return new memdb schema
// structs for constructing an in-memory db.
type schemaFn func() *memdb.TableSchema

// schemas is used to register schemas with the state store. Other files should
// use the init() callback to append to this.
var schemas []schemaFn

// stateStoreSchema is used to return the combined schema for the state store.
func stateStoreSchema() *memdb.DBSchema {
	// Create the root DB schema
	db := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema),
	}

	// Add the tables to the root schema
	for _, fn := range schemas {
		schema := fn()
		if _, ok := db.Tables[schema.Name]; ok {
			panic(fmt.Sprintf("duplicate table name: %s", schema.Name))
		}

		db.Tables[schema.Name] = schema
	}

	return db
}
