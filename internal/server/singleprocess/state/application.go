package state

import (
	"strings"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

// AppPut creates or updates the application.
func (s *State) AppPut(app *pb.Application) (*pb.Application, error) {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.appPut(dbTxn, memTxn, app)
	})
	if err == nil {
		memTxn.Commit()
	}

	return app, err
}

// AppDelete deletes an application from a project. This will also delete
// all the operations associated with this application.
func (s *State) AppDelete(ref *pb.Ref_Application) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.appDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// AppGet retrieves the application..
func (s *State) AppGet(ref *pb.Ref_Application) (*pb.Application, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Application
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.appGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

func (s *State) appPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Application,
) error {
	// Get the project
	p, err := s.projectGetOrCreate(dbTxn, memTxn, value.Project)
	if err != nil {
		return err
	}

	// If we have a matching app, then modify that that.
	pt := &serverptypes.Project{Project: p}
	if idx := pt.App(value.Name); idx >= 0 {
		p.Applications[idx] = value
		value = nil
	}

	// If we didn't have a matching app, insert it
	if value != nil {
		p.Applications = append(p.Applications, value)
	}

	return s.projectPut(dbTxn, memTxn, p)
}

func (s *State) appDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
) error {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: ref.Project,
	})
	if err != nil {
		// If the project doesn't exist then the app is deleted.
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// If we have a matching app, then modify that that.
	pt := &serverptypes.Project{Project: p}
	if i := pt.App(ref.Application); i >= 0 {
		s := p.Applications
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		p.Applications = s[:len(s)-1]
	}

	return nil
}

func (s *State) appGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
) (*pb.Application, error) {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: ref.Project,
	})
	if err != nil {
		return nil, err
	}

	// If we have a matching app, then modify that that.
	pt := &serverptypes.Project{Project: p}
	if i := pt.App(ref.Application); i >= 0 {
		return p.Applications[i], nil
	}

	return nil, status.Errorf(codes.NotFound, "application not found")
}

// appDefaultForRef returns a default pb.Application for a ref. This
// can be used in tandem with appCreateIfNotExist to create defaults.
func (s *State) appDefaultForRef(ref *pb.Ref_Application) *pb.Application {
	return &pb.Application{
		Name: ref.Application,
		Project: &pb.Ref_Project{
			Project: ref.Project,
		},
	}
}

// appIncrSeq gets the current sequence number, increments it and returns
// the one to set.
func (s *State) appIncrSeq(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
) (uint64, error) {
	// Get the project
	raw, err := memTxn.First(
		projectIndexTableName,
		projectIndexIdIndexName,
		string(s.projectIdByRef(&pb.Ref_Project{Project: ref.Project})),
	)
	if err != nil {
		return 0, err
	}

	ptr, ok := raw.(*projectIndexRecord).AppSeq[strings.ToLower(ref.Application)]
	if !ok {
		return 0, status.Errorf(codes.NotFound, "application not found")
	}

	return atomic.AddUint64(ptr, 1), nil
}

// appInitSeq sets the sequence number for an application if it is bigger.
func (s *State) appInitSeq(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
	value uint64,
) error {
	// Get the project
	raw, err := memTxn.First(
		projectIndexTableName,
		projectIndexIdIndexName,
		string(s.projectIdByRef(&pb.Ref_Project{Project: ref.Project})),
	)
	if err != nil {
		return err
	}

	ptr, ok := raw.(*projectIndexRecord).AppSeq[strings.ToLower(ref.Application)]
	if !ok {
		return status.Errorf(codes.NotFound, "application not found")
	}

	if value > *ptr {
		*ptr = value
	}

	return nil
}
