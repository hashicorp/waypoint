package singleprocess

import (
	"context"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var pushBucket = []byte("pushed_artifacts")

func init() {
	dbBuckets = append(dbBuckets, pushBucket)
}

func (s *service) UpsertPushedArtifact(
	ctx context.Context,
	req *pb.UpsertPushedArtifactRequest,
) (*pb.UpsertPushedArtifactResponse, error) {
	result := req.Artifact

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
		return dbUpsert(tx.Bucket(pushBucket), !insert, result.Id, result)
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpsertPushedArtifactResponse{Artifact: result}, nil
}

// TODO: test
func (s *service) ListPushedArtifacts(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListPushedArtifactsResponse, error) {
	var result []*pb.PushedArtifact
	s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(pushBucket)
		return bucket.ForEach(func(k, v []byte) error {
			var push pb.PushedArtifact
			if err := proto.Unmarshal(v, &push); err != nil {
				panic(err)
			}

			result = append(result, &push)
			return nil
		})
	})

	return &pb.ListPushedArtifactsResponse{Artifacts: result}, nil
}

// TODO: test
func (s *service) GetLatestPushedArtifact(
	ctx context.Context,
	req *empty.Empty,
) (*pb.PushedArtifact, error) {
	var result *pb.PushedArtifact
	var resultTime time.Time
	s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(pushBucket)
		return bucket.ForEach(func(k, v []byte) error {
			var push pb.PushedArtifact
			if err := proto.Unmarshal(v, &push); err != nil {
				panic(err)
			}

			// Looking for the push that is complete
			if push.Status.State != pb.Status_SUCCESS {
				return nil
			}

			t, err := ptypes.Timestamp(push.Status.CompleteTime)
			if err != nil {
				return status.Errorf(codes.Internal, "time for push can't be parsed")
			}

			if result == nil || resultTime.Before(t) {
				result = &push
				resultTime = t
			}

			return nil
		})
	})

	if result == nil {
		return nil, status.Errorf(codes.NotFound, "no successful pushes")
	}

	return result, nil
}
