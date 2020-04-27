package singleprocess

import (
	"context"

	"github.com/boltdb/bolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var releaseBucket = []byte("releases")

func init() {
	dbBuckets = append(dbBuckets, releaseBucket)
}

func (s *service) UpsertRelease(
	ctx context.Context,
	req *pb.UpsertReleaseRequest,
) (*pb.UpsertReleaseResponse, error) {
	result := req.Release

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
		return dbUpsert(tx.Bucket(releaseBucket), !insert, result.Id, result)
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpsertReleaseResponse{Release: result}, nil
}
