package state

import (
	"time"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
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

// AppPollPeek peeks at the next available project that will be polled against,
// and returns the projects application as the result along with the poll time.
// For more information on how ProjectPollPeek works, refer to the ProjectPollPeek
// docs.
func (s *State) AppPollPeek(
	ws memdb.WatchSet,
) (*pb.Application, time.Time, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Application
	var pollTime time.Time
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, pollTime, err = s.appPollPeek(dbTxn, memTxn, ws)
		return err
	})

	return result, pollTime, err
}

// AppPollComplete sets the next poll time for a given project given the app
// reference along with the time interval "t".
func (s *State) AppPollComplete(
	ref *pb.Ref_Application,
	t time.Time,
) error {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		err = s.appPollComplete(dbTxn, memTxn, ref, t)
		return err
	})

	return err
}

// GetFileChangeSignal checks the metadata for the given application and its
// project, returning the value of FileChangeSignal that is most relevent.
func (s *State) GetFileChangeSignal(scope *pb.Ref_Application) (string, error) {
	app, err := s.AppGet(scope)
	if err != nil {
		return "", err
	}

	if app.FileChangeSignal != "" {
		return app.FileChangeSignal, nil
	}

	project, err := s.ProjectGet(&pb.Ref_Project{Project: scope.Project})
	if err != nil {
		return "", err
	}

	return project.FileChangeSignal, nil
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

// TODO: instead of directly calling project poll peek, we probably will need
// to implement our own that peeks on the last applications poll time
func (s *State) appPollPeek(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
) (*pb.Application, time.Time, error) {
	// Get the project
	// Get next poll app from next project
	p, pollTime, err := s.ProjectPollPeek(ws)
	if err != nil {
		return nil, time.Time{}, err
	}

	// TODO: which application do we return in a project?
	return p.Applications[0], pollTime, nil
}

// TODO: if we have to implement our own peek, then we'll also want to add
// our down complete which completes the appl poll next time instead of project
func (s *State) appPollComplete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Application,
	t time.Time,
) error {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: ref.Project,
	})
	if err != nil {
		return err
	}

	err = s.ProjectPollComplete(p, t)
	if err != nil {
		return err
	}

	return nil
}
