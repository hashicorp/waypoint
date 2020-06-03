package singleprocess

import (
	"context"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var configBucket = []byte("config")

func init() {
	dbBuckets = append(dbBuckets, deployBucket)
}

func (s *service) SetConfig(
	ctx context.Context,
	req *pb.ConfigSetRequest,
) (*pb.ConfigSetResponse, error) {
	// Insert into our database
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(configBucket)
		if err != nil {
			return err
		}

		for _, cv := range req.Variables {
			err = dbPut(b, cv.Name, cv)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.ConfigSetResponse{}, nil
}

func (s *service) GetConfig(
	ctx context.Context,
	req *pb.ConfigGetRequest,
) (*pb.ConfigGetResponse, error) {
	var resp pb.ConfigGetResponse

	// Insert into our database
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(configBucket)
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			name := string(k)

			if req.Prefix == "" || strings.HasPrefix(name, req.Prefix) {
				var cv pb.ConfigVar

				err := proto.Unmarshal(v, &cv)
				if err != nil {
					return err
				}
				resp.Variables = append(resp.Variables, &cv)
			}

			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
