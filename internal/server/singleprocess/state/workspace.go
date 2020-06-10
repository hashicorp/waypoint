package state

import (
	"strings"

	"github.com/boltdb/bolt"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var workspaceBucket = []byte("workspace")

func init() {
	dbBuckets = append(dbBuckets, workspaceBucket)
}

// workspaceDefaultForRef returns a default pb.Workspace for a ref. This
// can be used in tandem with workspaceCreateIfNotExist to create defaults.
func (s *State) workspaceDefaultForRef(ref *pb.Ref_Workspace) *pb.Workspace {
	return &pb.Workspace{
		Name: ref.Workspace,
	}
}

func (s *State) workspaceCreateIfNotExist(tx *bolt.Tx, p *pb.Workspace) error {
	b, err := s.workspaceBucket(tx, p)
	if err != nil {
		return nil
	}

	id := []byte("value")
	if b.Get(id) != nil {
		// Workspace already exists
		return nil
	}

	// Create the workspace
	return dbPut(b, id, p)
}

func (s *State) workspaceBucket(tx *bolt.Tx, p *pb.Workspace) (*bolt.Bucket, error) {
	return tx.Bucket(workspaceBucket).CreateBucketIfNotExists(s.workspaceId(p))
}

func (s *State) workspaceId(p *pb.Workspace) []byte {
	return []byte(strings.ToLower(p.Name))
}
