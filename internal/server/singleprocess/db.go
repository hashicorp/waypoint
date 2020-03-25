package singleprocess

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// dbBuckets is the list of buckets that should be created by dbInit.
// Various components should use init() funcs to append to this.
var dbBuckets [][]byte

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

// dbUpsert is a helper to upsert a message. The update boolean will cause
// this to error if the ID is not found. This reflects our API behavior for
// upserts so that we don't let the end user pick any ID.
func dbUpsert(b *bolt.Bucket, update bool, id string, msg proto.Message) error {
	// If we're updating, the ID must exist
	if update && b.Get([]byte(id)) == nil {
		return status.Errorf(codes.NotFound, "record not found for ID: %s", id)
	}

	// Insert
	return dbPut(b, id, msg)
}

// dbPut is a helper to insert a proto.Message into a bucket for the given id.
// Any errors are automatically wrapped into a gRPC status error so they can
// be sent directly back.
func dbPut(b *bolt.Bucket, id string, msg proto.Message) error {
	enc, err := proto.Marshal(msg)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to encode data: %s", err)
	}

	if err := b.Put([]byte(id), enc); err != nil {
		return status.Errorf(codes.Aborted, "failed to write data: %s", err)
	}

	return nil
}

// dbGet is a helper to get a single proto.Message from a bucket. Errors
// are guaranteed to be in gRPC status format.
func dbGet(b *bolt.Bucket, id string, msg proto.Message) error {
	raw := b.Get([]byte(id))
	if raw == nil {
		return status.Errorf(codes.NotFound, "record not found for ID: %s", id)
	}

	if err := proto.Unmarshal(raw, msg); err != nil {
		return status.Errorf(codes.Internal, "failed to decode data: %s", err)
	}

	return nil
}

// dbList is a helper to list all the values in a bucket into a slice.
// The result should be a pointer to a typed slice of proto messages, example
// `[]*pb.Build`. This uses reflection to allocate the proper type and decode
// the messages.
func dbList(b *bolt.Bucket, result interface{}) error {
	return nil
}
