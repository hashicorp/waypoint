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

// TODO: test
func (s *service) ListDeployments(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListDeploymentsResponse, error) {
	var result []*pb.Deployment
	s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(deployBucket)
		return bucket.ForEach(func(k, v []byte) error {
			var deploy pb.Deployment
			if err := proto.Unmarshal(v, &deploy); err != nil {
				panic(err)
			}

			result = append(result, &deploy)
			return nil
		})
	})

	return &pb.ListDeploymentsResponse{Deployments: result}, nil
}
