package singleprocess

import (
	"context"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

var buildBucket = []byte("build")

func init() {
	dbBuckets = append(dbBuckets, buildBucket)
}

func (s *service) UpsertBuild(
	ctx context.Context,
	req *pb.UpsertBuildRequest,
) (*pb.UpsertBuildResponse, error) {
	result := req.Build

	// If we have no ID, then we're inserting and need to generate an ID.
	insert := result.Id == ""
	if insert {
		// Get the next id
		id, err := server.Id()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "uuid generation failed: %s", err)
		}

		// Specify the id
		result.Id = id
	}

	// Insert into our database
	err := s.db.Update(func(tx *bolt.Tx) error {
		return dbUpsert(tx.Bucket(buildBucket), !insert, result.Id, result)
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpsertBuildResponse{Build: result}, nil
}

func (s *service) ListBuilds(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListBuildsResponse, error) {
	var result []*pb.Build
	s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(buildBucket)
		return bucket.ForEach(func(k, v []byte) error {
			var build pb.Build
			if err := proto.Unmarshal(v, &build); err != nil {
				panic(err)
			}

			result = append(result, &build)
			return nil
		})
	})

	return &pb.ListBuildsResponse{Builds: result}, nil
}
