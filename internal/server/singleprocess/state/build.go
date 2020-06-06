package state

import (
	"github.com/boltdb/bolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var buildBucket = []byte("build")

func init() {
	dbBuckets = append(dbBuckets, buildBucket)
}

func (s *State) BuildPut(update bool, b *pb.Build) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return s.buildPut(tx, update, b)
	})
}

func (s *State) buildPut(tx *bolt.Tx, update bool, build *pb.Build) error {
	// Get our bucket
	b, err := s.buildBucket(tx, build)
	if err != nil {
		return err
	}

	// If this is an update, then the record should exist already
	id := []byte("value")
	if update && b.Get(id) == nil {
		return status.Errorf(codes.NotFound, "record not found for ID: %s", build.Id)
	}

	// Write our data
	return dbPut(b, id, build)
}

func (s *State) buildBucket(tx *bolt.Tx, build *pb.Build) (*bolt.Bucket, error) {
	b, err := s.appChildBucket(tx, buildBucket, &pb.Application{
		Name: build.Application.Application,
		Project: &pb.Ref_Project{
			Project: build.Application.Project,
		},
	})
	if err != nil {
		return nil, err
	}

	// Create the bucket for this build
	return b.CreateBucketIfNotExists([]byte(build.Id))
}
