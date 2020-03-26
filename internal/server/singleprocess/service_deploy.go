package singleprocess

import (
	"context"

	"github.com/boltdb/bolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
)

var deployBucket = []byte("deployments")

func init() {
	dbBuckets = append(dbBuckets, deployBucket)
}

func (s *service) UpsertDeployment(
	ctx context.Context,
	req *pb.UpsertDeploymentRequest,
) (*pb.UpsertDeploymentResponse, error) {
	result := req.Deployment

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
		return dbUpsert(tx.Bucket(deployBucket), !insert, result.Id, result)
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpsertDeploymentResponse{Deployment: result}, nil
}
