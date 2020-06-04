package singleprocess

import (
	"context"
	"sort"
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
			if cv.Value == "" {
				err = b.Delete([]byte(cv.App + ":" + cv.Name))
			} else {
				err = dbPut(b, cv.App+":"+cv.Name, cv)
			}

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

		if req.App != "" {
			vars := map[string]*pb.ConfigVar{}
			var keys []string

			err := b.ForEach(func(k, v []byte) error {
				name := string(k)

				if req.Prefix == "" || strings.HasPrefix(name, req.Prefix) {
					var cv pb.ConfigVar

					err := proto.Unmarshal(v, &cv)
					if err != nil {
						return err
					}

					if cv.App != req.App {
						return nil
					}

					cur := vars[cv.Name]
					if cur != nil && cur.App != "" {
						return nil
					}

					vars[cv.Name] = &cv
					keys = append(keys, cv.Name)
					return nil
				}

				return nil
			})
			if err != nil {
				return err
			}

			sort.Strings(keys)

			for _, k := range keys {
				resp.Variables = append(resp.Variables, vars[k])
			}
		} else {
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
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
