// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"strings"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var authMethodBucket = []byte("auth_method")

func init() {
	dbBuckets = append(dbBuckets, authMethodBucket)
	dbIndexers = append(dbIndexers, (*State).authMethodIndexInit)
	schemas = append(schemas, authMethodSchema)
}

// AuthMethodPut creates or updates the given auth method. .
func (s *State) AuthMethodPut(ctx context.Context, v *pb.AuthMethod) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.authMethodPut(dbTxn, memTxn, v)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// AuthMethodGet gets a auth method by reference.
func (s *State) AuthMethodGet(ctx context.Context, ref *pb.Ref_AuthMethod) (*pb.AuthMethod, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.AuthMethod
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.authMethodGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// AuthMethodDelete deletes an auth method by reference.
func (s *State) AuthMethodDelete(ctx context.Context, ref *pb.Ref_AuthMethod) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.authMethodDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// AuthMethodList returns the list of projects.
func (s *State) AuthMethodList(ctx context.Context) ([]*pb.AuthMethod, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(authMethodTableName, authMethodIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.AuthMethod
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*authMethodIndex)

		var v *pb.AuthMethod
		err = s.db.View(func(dbTxn *bolt.Tx) error {
			v, err = s.authMethodGet(dbTxn, memTxn, &pb.Ref_AuthMethod{Name: idx.Id})
			return err
		})

		result = append(result, v)
	}

	return result, nil
}

func (s *State) authMethodPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.AuthMethod,
) error {
	id := s.authMethodId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(authMethodBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.authMethodIndexSet(memTxn, id, value)
}

func (s *State) authMethodGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_AuthMethod,
) (*pb.AuthMethod, error) {
	var result pb.AuthMethod
	b := dbTxn.Bucket(authMethodBucket)
	return &result, dbGet(b, []byte(strings.ToLower(ref.Name)), &result)
}

func (s *State) authMethodDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_AuthMethod,
) error {
	// Get the authMethod. If it doesn't exist then we're successful.
	v, err := s.authMethodGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete from bolt
	id := s.authMethodId(v)
	bucket := dbTxn.Bucket(authMethodBucket)
	if err := bucket.Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(authMethodTableName, &authMethodIndex{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// authMethodIndexSet writes an index record for a single authMethod.
func (s *State) authMethodIndexSet(txn *memdb.Txn, id []byte, value *pb.AuthMethod) error {
	record := &authMethodIndex{
		Id: string(id),
	}

	// Insert the index
	return txn.Insert(authMethodTableName, record)
}

// authMethodIndexInit initializes the authMethod index from persisted data.
func (s *State) authMethodIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(authMethodBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.AuthMethod
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.authMethodIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) authMethodId(v *pb.AuthMethod) []byte {
	return []byte(strings.ToLower(v.Name))
}

func authMethodSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: authMethodTableName,
		Indexes: map[string]*memdb.IndexSchema{
			authMethodIdIndexName: {
				Name:         authMethodIdIndexName,
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
	authMethodTableName   = "auth-methods"
	authMethodIdIndexName = "id"
)

type authMethodIndex struct {
	Id string // Unique ID, we just use the auth method name
}

// Copy should be called prior to any modifications to an existing record.
func (idx *authMethodIndex) Copy() *authMethodIndex {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
