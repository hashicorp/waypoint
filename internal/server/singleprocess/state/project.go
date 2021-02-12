package state

import (
	"strings"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var projectBucket = []byte("project")

func init() {
	dbBuckets = append(dbBuckets, projectBucket)
	dbIndexers = append(dbIndexers, (*State).projectIndexInit)
	schemas = append(schemas, projectIndexSchema)
}

// ProjectPut creates or updates the given project.
//
// Application changes will be ignored, you must use the Application APIs.
func (s *State) ProjectPut(p *pb.Project) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		prev, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
			Project: p.Name,
		})
		if err != nil && status.Code(err) != codes.NotFound {
			// We ignore NotFound since this function is used to create projects.
			return err
		}
		if err == nil {
			// If we have a previous project, preserve the applications.
			p.Applications = prev.Applications
		}

		return s.projectPut(dbTxn, memTxn, p)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// ProjectGet gets a project by reference.
func (s *State) ProjectGet(ref *pb.Ref_Project) (*pb.Project, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Project
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.projectGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// ProjectDelete deletes a project by reference. This is a complete data
// delete. This will delete all operations associated with this project
// as well.
func (s *State) ProjectDelete(ref *pb.Ref_Project) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.projectDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// ProjectList returns the list of projects.
func (s *State) ProjectList() ([]*pb.Ref_Project, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.projectList(memTxn)
}

// ProjectUpdateDataRef updates the latest data ref used for a project.
// This data is available via the APIs for querying workspaces.
func (s *State) ProjectUpdateDataRef(
	ref *pb.Ref_Project,
	ws *pb.Ref_Workspace,
	dataRef *pb.Job_DataSource_Ref,
) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.workspaceUpdateProjectDataRef(dbTxn, memTxn, ws, ref, dataRef)
	})
	if err != nil {
		return err
	}

	memTxn.Commit()
	return nil
}

func (s *State) projectGetOrCreate(dbTxn *bolt.Tx, memTxn *memdb.Txn, ref *pb.Ref_Project) (*pb.Project, error) {
	result, err := s.projectGet(dbTxn, memTxn, ref)
	if status.Code(err) == codes.NotFound {
		result = nil
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	result = &pb.Project{
		Name: ref.Project,
	}

	return result, s.projectPut(dbTxn, memTxn, result)
}

func (s *State) projectPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Project,
) error {
	// This is to prevent mistakes or abuse. Realistically a waypoint.hcl
	// file should be MUCH smaller than this so this catches the really big
	// mistakes.
	if len(value.WaypointHcl) > projectWaypointHclMaxSize {
		return status.Errorf(codes.FailedPrecondition,
			"project 'waypoint_hcl' exceeds maximum size (5MB)",
		)
	}

	id := s.projectId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(projectBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.projectIndexSet(memTxn, id, value)
}

func (s *State) projectGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Project,
) (*pb.Project, error) {
	var result pb.Project
	b := dbTxn.Bucket(projectBucket)
	return &result, dbGet(b, []byte(strings.ToLower(ref.Project)), &result)
}

func (s *State) projectList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_Project, error) {
	iter, err := memTxn.Get(projectIndexTableName, projectIndexIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_Project
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*projectIndexRecord)

		result = append(result, &pb.Ref_Project{
			Project: idx.Id,
		})
	}

	return result, nil
}

func (s *State) projectDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Project,
) error {
	// Get the project. If it doesn't exist then we're successful.
	p, err := s.projectGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete each application
	for _, app := range p.Applications {
		if err := s.appDelete(dbTxn, memTxn, &pb.Ref_Application{
			Project:     ref.Project,
			Application: app.Name,
		}); err != nil {
			return err
		}
	}

	// Delete from bolt
	id := s.projectIdByRef(ref)
	if err := dbTxn.Bucket(projectBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(projectIndexTableName, &projectIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// projectIndexSet writes an index record for a single project.
func (s *State) projectIndexSet(txn *memdb.Txn, id []byte, value *pb.Project) error {
	record := &projectIndexRecord{
		Id: string(id),
	}

	// Insert the index
	return txn.Insert(projectIndexTableName, record)
}

// projectIndexInit initializes the project index from persisted data.
func (s *State) projectIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(projectBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.Project
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.projectIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) projectId(p *pb.Project) []byte {
	return []byte(strings.ToLower(p.Name))
}

func (s *State) projectIdByRef(ref *pb.Ref_Project) []byte {
	return []byte(strings.ToLower(ref.Project))
}

func projectIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: projectIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			projectIndexIdIndexName: {
				Name:         projectIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	projectIndexTableName   = "project-index"
	projectIndexIdIndexName = "id"

	projectWaypointHclMaxSize = 5 * 1024 // 5 MB
)

type projectIndexRecord struct {
	Id string
}
