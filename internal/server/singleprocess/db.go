package singleprocess

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
