package state

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var userBucket = []byte("user")

func init() {
	dbBuckets = append(dbBuckets, userBucket)
	dbIndexers = append(dbIndexers, (*State).userIndexInit)
	schemas = append(schemas, userIndexSchema)
}

// UserPut creates or updates the given user. If the user has no ID set
// then an ID will be written directly to the parameter.
func (s *State) UserPut(user *pb.User) error {
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
func (s *State) UserGet(ref *pb.Ref_User) (*pb.User, error) {
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
func (s *State) UserDelete(ref *pb.Ref_User) error {
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
func (s *State) UserList() ([]*pb.User, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(userIndexTableName, userIndexIdIndexName+"_prefix", "")
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
func (s *State) UserEmpty() (bool, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(userIndexTableName, userIndexIdIndexName+"_prefix", "")
	if err != nil {
		return false, err
	}

	return iter.Next() == nil, nil
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
			userIndexTableName,
			userIndexUsernameIndexName,
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
	if err := memTxn.Delete(userIndexTableName, &userIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// userIndexSet writes an index record for a single user.
func (s *State) userIndexSet(txn *memdb.Txn, id []byte, value *pb.User) error {
	record := &userIndexRecord{
		Id:       string(id),
		Username: value.Username,
	}

	// Insert the index
	return txn.Insert(userIndexTableName, record)
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
		Name: userIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			userIndexIdIndexName: {
				Name:         userIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			userIndexUsernameIndexName: {
				Name:         userIndexUsernameIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Username",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	userIndexTableName         = "user-index"
	userIndexIdIndexName       = "id"
	userIndexUsernameIndexName = "username"
)

type userIndexRecord struct {
	Id       string
	Username string
}

// Copy should be called prior to any modifications to an existing record.
func (idx *userIndexRecord) Copy() *userIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
