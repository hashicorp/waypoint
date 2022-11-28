package boltdbstate

import (
	"context"
	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var (
	serverConfigBucket = []byte("server-config")
	serverConfigId     = []byte("1")
)

func init() {
	dbBuckets = append(dbBuckets, serverConfigBucket)
	dbIndexers = append(dbIndexers, (*State).serverConfigIndexInit)
	schemas = append(schemas, serverConfigIndexSchema)
}

// ServerConfigSet writes the server configuration.
func (s *State) ServerConfigSet(ctx context.Context, c *pb.ServerConfig) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.serverConfigSet(dbTxn, memTxn, c)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// ServerConfigGet gets the server configuration.
func (s *State) ServerConfigGet(ctx context.Context) (*pb.ServerConfig, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.serverConfigGet(memTxn)
}

func (s *State) serverConfigGet(
	memTxn *memdb.Txn,
) (*pb.ServerConfig, error) {
	v, err := memTxn.First(
		serverConfigIndexTableName,
		serverConfigIndexIdIndexName,
		string(serverConfigId),
	)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return &pb.ServerConfig{}, nil
	}

	return v.(*serverConfigIndexRecord).Config, nil
}

func (s *State) serverConfigSet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.ServerConfig,
) error {
	id := serverConfigId

	// If we want to delete the config, just set it to empty so we can keep
	// the server-set settings such as cookies.
	if value == nil {
		value = &pb.ServerConfig{}
	}

	// Get the previous value. We do not allow overwriting the cookie and we
	// also will initialize the cookie if it is not set.
	old, err := s.serverConfigGet(memTxn)
	if err != nil {
		return err
	}
	value.Cookie = ""
	if old != nil {
		value.Cookie = old.Cookie
	}

	// Generate a cookie if we don't have one
	if value.Cookie == "" {
		newCookie, err := ulid()
		if err != nil {
			return err
		}

		value.Cookie = "WPC" + newCookie
	}

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(serverConfigBucket)
	if value == nil {
		if err := b.Delete(id); err != nil {
			return err
		}
	} else {
		if err := dbPut(b, id, value); err != nil {
			return err
		}
	}

	// Create our index value and write that.
	return s.serverConfigIndexSet(memTxn, id, value)
}

// serverConfigIndexSet writes an index record for the server config.
func (s *State) serverConfigIndexSet(txn *memdb.Txn, id []byte, value *pb.ServerConfig) error {
	record := &serverConfigIndexRecord{
		Id:     string(id),
		Config: value,
	}

	// If we have no value, we delete from the memdb index
	if value == nil {
		return txn.Delete(serverConfigIndexTableName, record)
	}

	// Insert the index
	return txn.Insert(serverConfigIndexTableName, record)
}

// serverConfigIndexInit initializes the server config index from persisted data.
func (s *State) serverConfigIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(serverConfigBucket)

	data := bucket.Get(serverConfigId)
	if data == nil {
		return nil
	}

	var value pb.ServerConfig
	if err := proto.Unmarshal(data, &value); err != nil {
		return err
	}
	if err := s.serverConfigIndexSet(memTxn, serverConfigId, &value); err != nil {
		return err
	}

	return nil
}

func serverConfigIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: serverConfigIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			serverConfigIndexIdIndexName: {
				Name:         serverConfigIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	serverConfigIndexTableName   = "server-config-index"
	serverConfigIndexIdIndexName = "id"
)

// Our record for the server config index. We will only have at most one
// of these because server config is a singleton.
type serverConfigIndexRecord struct {
	Id     string
	Config *pb.ServerConfig
}
