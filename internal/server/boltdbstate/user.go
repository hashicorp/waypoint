// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"context"
	"strings"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/ptypes"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

var userBucket = []byte("user")

func init() {
	dbBuckets = append(dbBuckets, userBucket)
	dbIndexers = append(dbIndexers, (*State).userIndexInit)
	schemas = append(schemas, userIndexSchema)
}

// UserPut creates or updates the given user. If the user has no ID set
// then an ID will be written directly to the parameter.
func (s *State) UserPut(ctx context.Context, user *pb.User) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.userPut(dbTxn, memTxn, user)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// UserGet gets a user by reference.
func (s *State) UserGet(ctx context.Context, ref *pb.Ref_User) (*pb.User, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.User
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.userGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// UserDelete deletes a user by reference.
func (s *State) UserDelete(ctx context.Context, ref *pb.Ref_User) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.userDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// UserList returns the list of projects.
func (s *State) UserList(ctx context.Context) ([]*pb.User, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(userTableName, userIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.User
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*userIndexRecord)

		var v *pb.User
		err = s.db.View(func(dbTxn *bolt.Tx) error {
			v, err = s.userGet(dbTxn, memTxn, &pb.Ref_User{
				Ref: &pb.Ref_User_Id{
					Id: &pb.Ref_UserId{Id: idx.Id},
				},
			})
			return err
		})

		result = append(result, v)
	}

	return result, nil
}

// UserEmpty returns true if there are no users yet (bootstrap state).
func (s *State) UserEmpty(ctx context.Context) (bool, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(userTableName, userIdIndexName+"_prefix", "")
	if err != nil {
		return false, err
	}

	return iter.Next() == nil, nil
}

// UserGetOIDC gets a user by by OIDC link lookup. This will return
// a codes.NotFound error if the user is not found.
func (s *State) UserGetOIDC(ctx context.Context, iss, sub string) (*pb.User, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.User
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.userGetOIDC(dbTxn, memTxn, iss, sub)
		return err
	})

	return result, err
}

func (s *State) userGetOIDC(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	iss, sub string,
) (*pb.User, error) {
	b := dbTxn.Bucket(userBucket)

	// Look up all users that match this sub.
	iter, err := memTxn.Get(
		userTableName,
		userOIDCSubIndexName,
		sub,
	)
	if err != nil {
		return nil, err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		idx := raw.(*userIndexRecord)

		// Read the user from disk
		var result pb.User
		if err := dbGet(b, []byte(strings.ToLower(idx.Id)), &result); err != nil {
			return nil, err
		}

		// Compare the issuer. We need to find a user with a single link that
		// has both the issuer and sub.
		for _, link := range result.Links {
			oidc, ok := link.Method.(*pb.User_Link_Oidc)
			if !ok || oidc == nil {
				continue
			}

			if strings.EqualFold(oidc.Oidc.Sub, sub) && strings.EqualFold(oidc.Oidc.Iss, iss) {
				return &result, nil
			}
		}
	}

	return nil, status.Errorf(codes.NotFound, "user not found")
}

// UserGetEmail gets a user by by email lookup. This will return
// a codes.NotFound error if the user is not found.
func (s *State) UserGetEmail(ctx context.Context, email string) (*pb.User, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.User
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.userGetEmail(dbTxn, memTxn, email)
		return err
	})

	return result, err
}

func (s *State) userGetEmail(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	email string,
) (*pb.User, error) {
	b := dbTxn.Bucket(userBucket)

	// Look up all users that match this sub.
	iter, err := memTxn.Get(
		userTableName,
		userEmailIndexName,
		email,
	)
	if err != nil {
		return nil, err
	}

	raw := iter.Next()
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	idx := raw.(*userIndexRecord)

	var result pb.User
	return &result, dbGet(b, []byte(strings.ToLower(idx.Id)), &result)
}

func (s *State) userPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.User,
) error {
	// If the user doesn't have an ID set, we create one.
	if value.Id == "" {
		id, err := ulid()
		if err != nil {
			return err
		}

		value.Id = id
	}

	// We want to validate the user again.
	if err := ptypes.ValidateUser(value); err != nil {
		return err
	}

	id := s.userId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(userBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.userIndexSet(memTxn, id, value)
}

func (s *State) userGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_User,
) (*pb.User, error) {
	var result pb.User
	b := dbTxn.Bucket(userBucket)

	switch ref := ref.Ref.(type) {
	case *pb.Ref_User_Id:
		return &result, dbGet(b, []byte(strings.ToLower(ref.Id.Id)), &result)

	case *pb.Ref_User_Username:
		iter, err := memTxn.Get(
			userTableName,
			userUsernameIndexName,
			ref.Username.Username,
		)
		if err != nil {
			return nil, err
		}

		raw := iter.Next()
		if raw == nil {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		idx := raw.(*userIndexRecord)

		return &result, dbGet(b, []byte(strings.ToLower(idx.Id)), &result)

	default:
		return nil, status.Errorf(codes.FailedPrecondition, "invalid user ref type")
	}
}

func (s *State) userDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_User,
) error {
	// Get the user. If it doesn't exist then we're successful.
	u, err := s.userGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// If the user is the default user, then we can't delete them for now
	if u.Id == serverstate.DefaultUserId {
		return status.Errorf(codes.FailedPrecondition,
			"The initial Waypoint user can't currently be deleted. The initial "+
				"user is used by deployments and runners for authentication. "+
				"A future version of Waypoint will remove this restriction.")
	}

	// We can't delete the final user or the system will get into a state
	// where it can't do anything!
	bucket := dbTxn.Bucket(userBucket)
	if bucket.Stats().KeyN <= 1 {
		return status.Errorf(codes.FailedPrecondition,
			"the final user in the system can't be deleted until another one is created")
	}

	// Delete from bolt
	id := s.userId(u)
	if err := bucket.Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(userTableName, &userIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// userIndexSet writes an index record for a single user.
func (s *State) userIndexSet(txn *memdb.Txn, id []byte, value *pb.User) error {
	record := &userIndexRecord{
		Id:       string(id),
		Username: value.Username,
		Email:    value.Email,
	}

	for _, link := range value.Links {
		switch method := link.Method.(type) {
		case *pb.User_Link_Oidc:
			record.OIDCSub = append(record.OIDCSub, method.Oidc.Sub)
		}
	}

	// Insert the index
	return txn.Insert(userTableName, record)
}

// userIndexInit initializes the user index from persisted data.
func (s *State) userIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(userBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.User
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.userIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) userId(u *pb.User) []byte {
	return []byte(strings.ToLower(u.Id))
}

func userIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: userTableName,
		Indexes: map[string]*memdb.IndexSchema{
			userIdIndexName: {
				Name:         userIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			userUsernameIndexName: {
				Name:         userUsernameIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Username",
					Lowercase: true,
				},
			},

			userEmailIndexName: {
				Name:         userEmailIndexName,
				AllowMissing: true,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Email",
					Lowercase: true,
				},
			},

			userOIDCSubIndexName: {
				Name:         userOIDCSubIndexName,
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field:     "OIDCSub",
					Lowercase: true,
				},

				// This field is almost always unique but isn't guaranteed
				// since uniqueness depends on the issuer + sub combo.
				Unique: false,
			},
		},
	}
}

const (
	userTableName         = "user-index"
	userIdIndexName       = "id"
	userUsernameIndexName = "username"
	userEmailIndexName    = "email"
	userOIDCSubIndexName  = "oidc-sub"
)

type userIndexRecord struct {
	Id       string
	Username string
	Email    string

	// OIDCSub is a list of OIDC sub claims that are linked to this user.
	// This can be used to look up a user by OIDC.
	OIDCSub []string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *userIndexRecord) Copy() *userIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
