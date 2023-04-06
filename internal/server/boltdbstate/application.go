// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"time"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

// AppPut creates or updates the application.
func (s *State) AppPut(ctx context.Context, app *pb.Application) (*pb.Application, error) {
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
func (s *State) AppDelete(ctx context.Context, ref *pb.Ref_Application) error {
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
func (s *State) AppGet(ctx context.Context, ref *pb.Ref_Application) (*pb.Application, error) {
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

// ApplicationPollPeek peeks at the next available project that will be polled against,
// and returns the project as the result along with the poll time. The poll queuer
// will queue a job against every defined application for the given project.
// For more information on how ProjectPollPeek works, refer to the ProjectPollPeek
// docs.
func (s *State) ApplicationPollPeek(
	ctx context.Context,
	ws memdb.WatchSet,
) (*pb.Project, time.Time, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Project
	var pollTime time.Time
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, pollTime, err = s.appPollPeek(dbTxn, memTxn, ws)
		return err
	})

	return result, pollTime, err
}

// ApplicationPollComplete sets the next poll time for a given project given the app
// reference along with the time interval "t".
func (s *State) ApplicationPollComplete(
	ctx context.Context,
	project *pb.Project,
	t time.Time,
) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.View(func(dbTxn *bolt.Tx) error {
		err := s.appPollComplete(dbTxn, memTxn, project, t)
		return err
	})

	return err
}

// GetFileChangeSignal checks the metadata for the given application and its
// project, returning the value of FileChangeSignal that is most relevent.
func (s *State) GetFileChangeSignal(ctx context.Context, scope *pb.Ref_Application) (string, error) {
	app, err := s.AppGet(ctx, scope)
	if err != nil {
		return "", err
	}

	if app.FileChangeSignal != "" {
		return app.FileChangeSignal, nil
	}

	project, err := s.ProjectGet(ctx, &pb.Ref_Project{Project: scope.Project})
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

	// If we have a matching app, then modify that.
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

	// If we have a matching app, then modify that.
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

func (s *State) appPollPeek(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ws memdb.WatchSet,
) (*pb.Project, time.Time, error) {
	// LowerBound doesn't support watches so we have to do a Get first
	// to get a valid watch channel on these fields.
	iter, err := memTxn.Get(
		projectIndexTableName,
		appIndexNextPollIndexName,
		true,            // polling enabled
		time.Unix(0, 0), // lowest next poll time
	)
	if err != nil {
		return nil, time.Time{}, err
	}
	ws.Add(iter.WatchCh())

	// Get the project's app with the lowest "next poll" time.
	iter, err = memTxn.LowerBound(
		projectIndexTableName,
		appIndexNextPollIndexName,
		true,            // polling enabled
		time.Unix(0, 0), // lowest next poll time
	)
	if err != nil {
		return nil, time.Time{}, err
	}

	// If we have no values, then return
	raw := iter.Next()
	if raw == nil {
		return nil, time.Time{}, nil
	}

	rec := raw.(*projectIndexRecord)
	if rec.AppStatusNextPoll.IsZero() {
		// This happens if this application's poller hasn't been switched on
		return nil, time.Time{}, nil
	}

	var result *pb.Project

	result, err = s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: rec.Id,
	})

	return result, rec.AppStatusNextPoll, err
}

func (s *State) appPollComplete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	project *pb.Project,
	t time.Time,
) error {
	// Get the project
	p, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
		Project: project.Name,
	})
	if status.Code(err) == codes.NotFound {
		return nil
	}
	if err != nil {
		return err
	}

	raw, err := memTxn.First(
		projectIndexTableName,
		projectIndexIdIndexName,
		string(s.projectId(p)),
	)
	if err != nil {
		return err
	}
	if raw == nil {
		return nil
	}

	record := raw.(*projectIndexRecord)
	if !record.AppStatusPoll {
		// If this project doesn't have polling enabled, then do nothing.
		// This could happen if a project had polling when Peek was called,
		// then between Peek and Complete, polling was disabled.
		return nil
	}

	record = record.Copy()
	record.AppStatusLastPoll = t
	record.AppStatusNextPoll = t.Add(record.AppStatusPollInterval)
	if err := memTxn.Insert(projectIndexTableName, record); err != nil {
		return err
	}

	memTxn.Commit()
	return nil
}
