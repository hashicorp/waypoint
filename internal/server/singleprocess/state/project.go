package state

import (
	"github.com/boltdb/bolt"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var projectBucket = []byte("project")

func init() {
	dbBuckets = append(dbBuckets, projectBucket)
}

// projectDefaultForRef returns a default pb.Project for a ref. This
// can be used in tandem with projectCreateIfNotExist to create defaults.
func (s *State) projectDefaultForRef(ref *pb.Ref_Project) *pb.Project {
	return &pb.Project{
		Name: ref.Project,
	}
}

func (s *State) projectCreateIfNotExist(tx *bolt.Tx, p *pb.Project) error {
	b, err := s.projectBucket(tx, p)
	if err != nil {
		return nil
	}

	id := []byte("value")
	if b.Get(id) != nil {
		// Project already exists
		return nil
	}

	// Create the project
	return dbPut(b, id, p)
}

func (s *State) projectBucket(tx *bolt.Tx, p *pb.Project) (*bolt.Bucket, error) {
	return tx.Bucket(projectBucket).CreateBucketIfNotExists(s.projectId(p))
}

func (s *State) projectId(p *pb.Project) []byte {
	return []byte(p.Name)
}
