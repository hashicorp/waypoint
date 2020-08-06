package state

import (
	"crypto/rand"
	"io"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var hmacKeyOp = &appOperation{
	Struct: (*pb.HMACKey)(nil),
	Bucket: []byte("hmackey"),
}

func init() {
	hmacKeyOp.register()
}

// HMACKeyPut inserts or updates a build record.
func (s *State) HMACKeyCreate(id string, size int) (*pb.HMACKey, error) {
	var key pb.HMACKey
	key.Id = id

	err := s.db.Update(func(tx *bolt.Tx) error {
		// Get the global bucket and write the value to it.
		b := tx.Bucket(hmacKeyOp.Bucket)

		// If we're updating, then this shouldn't already exist
		bid := []byte(id)

		cur := b.Get(bid)
		if cur != nil {
			return proto.Unmarshal(cur, &key)
		}

		raw := make([]byte, size)

		_, err := io.ReadFull(rand.Reader, raw)
		if err != nil {
			return err
		}

		key.Id = id
		key.Key = raw

		return dbPut(b, bid, &key)
	})
	if err != nil {
		return nil, err
	}

	return &key, nil
}

// HMACKeyGet gets a build by ID.
func (s *State) HMACKeyGet(id string) (*pb.HMACKey, error) {
	result, err := hmacKeyOp.Get(s, &pb.Ref_Operation{
		Target: &pb.Ref_Operation_Id{Id: id},
	})
	if err != nil {
		return nil, err
	}

	return result.(*pb.HMACKey), nil
}
