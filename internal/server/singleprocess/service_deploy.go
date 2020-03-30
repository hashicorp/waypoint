package singleprocess

import (
	"context"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mitchellh/devflow/internal/server"
	pb "github.com/mitchellh/devflow/internal/server/gen"
	serversort "github.com/mitchellh/devflow/internal/server/sort"
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
	req *pb.ListDeploymentsRequest,
) (*pb.ListDeploymentsResponse, error) {
	var result []*pb.Deployment
	s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(deployBucket)
		return bucket.ForEach(func(k, v []byte) error {
			var deploy pb.Deployment
			if err := proto.Unmarshal(v, &deploy); err != nil {
				panic(err)
			}

			// Filter
			if !statusFilterMatch(req.Status, deploy.Status) {
				return nil
			}

			result = append(result, &deploy)
			return nil
		})
	})

	// Sort if we have to
	var sortIface sort.Interface
	switch req.Order {
	case pb.ListDeploymentsRequest_START_TIME:
		sortIface = serversort.DeploymentStartDesc(result)
		if !req.OrderDesc {
			sortIface = sort.Reverse(sortIface)
		}

	case pb.ListDeploymentsRequest_COMPLETE_TIME:
		sortIface = serversort.DeploymentCompleteDesc(result)
		if !req.OrderDesc {
			sortIface = sort.Reverse(sortIface)
		}
	}
	if sortIface != nil {
		sort.Sort(sortIface)
	}

	// Limit
	if req.Limit > 0 && req.Limit < uint32(len(result)) {
		result = result[:req.Limit]
	}

	return &pb.ListDeploymentsResponse{Deployments: result}, nil
}
