package state

import (
	"github.com/boltdb/bolt"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var appBucket = []byte("app")

func init() {
	dbBuckets = append(dbBuckets, appBucket)
}

func (s *State) appCreateIfNotExist(tx *bolt.Tx, app *pb.Application) error {
	// Create our project if we don't have it already.
	err := s.projectCreateIfNotExist(tx, &pb.Project{
		Name: app.Project.Project,
	})
	if err != nil {
		return err
	}

	// Get our bucket
	b, err := s.appBucket(tx, app)
	if err != nil {
		return err
	}

	id := []byte("value")
	if b.Get(id) != nil {
		return nil
	}

	// Write our data
	return dbPut(b, id, app)
}

func (s *State) appBucket(tx *bolt.Tx, app *pb.Application) (*bolt.Bucket, error) {
	// Get the app bucket within the specific project.
	appBucket, err := s.projectAppBucket(tx, &pb.Project{
		Name: app.Project.Project,
	})
	if err != nil {
		return nil, err
	}

	return appBucket.CreateBucketIfNotExists(s.appId(app))
}

func (s *State) appChildBucket(tx *bolt.Tx, name []byte, app *pb.Application) (*bolt.Bucket, error) {
	// Create our app if we don't have it already
	if err := s.appCreateIfNotExist(tx, app); err != nil {
		return nil, err
	}

	// Get the app bucket
	appBucket, err := s.appBucket(tx, app)
	if err != nil {
		return nil, err
	}

	// Create and return the child bucket
	return appBucket.CreateBucketIfNotExists(name)
}

func (s *State) appId(app *pb.Application) []byte {
	return []byte(app.Name)
}

func (s *State) projectAppBucket(tx *bolt.Tx, p *pb.Project) (*bolt.Bucket, error) {
	// Get the project bucket since applications live off of that.
	projBucket, err := s.projectBucket(tx, p)
	if err != nil {
		return nil, err
	}

	// Create the applications bucket
	return projBucket.CreateBucketIfNotExists(appBucket)
}
