// Package state manages the state that the singleprocess server has, providing
// operations to mutate that state safely as needed.
package state

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The global variables below can be set by init() functions of other
// files in this package to setup the database state for the server.
var (
	// schemas is used to register schemas with the state store. Other files should
	// use the init() callback to append to this.
	schemas []schemaFn

	// dbBuckets is the list of buckets that should be created by dbInit.
	// Various components should use init() funcs to append to this.
	dbBuckets [][]byte

	// dbIndexers is the list of functions to call to initialize the
	// in-memory indexes from the persisted db.
	dbIndexers []indexFn
)

// State is the primary API for state mutation for the server.
type State struct {
	// inmem is our in-memory database that stores ephemeral data in an
	// easier-to-query way. Some of this data may be periodically persisted
	// but most of this data is meant to be lost when the process restarts.
	inmem *memdb.MemDB

	// db is our persisted on-disk database. This stores the bulk of data
	// and supports a transactional model for safe concurrent access.
	// inmem is used alongside db to store in-memory indexing information
	// for more efficient lookups into db. This index is built online at
	// boot.
	db *bolt.DB
}

// New initializes a new State store.
func New(db *bolt.DB) (*State, error) {
	// Create the in-memory DB.
	inmem, err := memdb.NewMemDB(stateStoreSchema())
	if err != nil {
		return nil, fmt.Errorf("Failed setting up state store: %s", err)
	}

	// Initialize and validate our on-disk format.
	if err := dbInit(db); err != nil {
		return nil, err
	}

	s := &State{inmem: inmem, db: db}

	// Initialize our in-memory indexes
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()
	err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, indexer := range dbIndexers {
			if err := indexer(s, dbTxn, memTxn); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	memTxn.Commit()

	return s, nil
}

// Close should be called to gracefully close any resources.
func (s *State) Close() error {
	// Nothing for now, but we expect to do things one day.
	return nil
}

// schemaFn is an interface function used to create and return new memdb schema
// structs for constructing an in-memory db.
type schemaFn func() *memdb.TableSchema

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

// indexFn is the function type for initializing in-memory indexes from
// persisted data. This is usually specified as a method handle to a
// *State method.
//
// The bolt.Tx is read-only while the memdb.Txn is a write transaction.
type indexFn func(*State, *bolt.Tx, *memdb.Txn) error

// dbInit sets up the database. This should be called once on all new
// DB handles before accepting API calls. It is safe to be called multiple
// times.
func dbInit(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Create all our buckets
		for _, b := range dbBuckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}

		// Check our data version
		// TODO(mitchellh): make this work
		sys := tx.Bucket(sysBucket)
		vsnRaw := sys.Get(sysVersionKey)
		if len(vsnRaw) > 0 {
			return status.Errorf(
				codes.FailedPrecondition,
				"system version is set, shouldn't be yet",
			)
		}

		return nil
	})
}

var (
	// sysBucket stores system-related information.
	sysBucket = []byte("system")

	// sysVersionKey stores the version of the data that is stored.
	// This is used for data migration.
	sysVersionKey = []byte("version")
)

func init() {
	dbBuckets = append(dbBuckets, sysBucket)
}
