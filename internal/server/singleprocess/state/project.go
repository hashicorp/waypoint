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
func (s *State) ProjectPut(p *pb.Project) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
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
			projectIndexIdIndexName: &memdb.IndexSchema{
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
)

type projectIndexRecord struct {
	Id string
}
