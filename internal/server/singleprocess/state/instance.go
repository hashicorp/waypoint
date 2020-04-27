package state

import (
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server/logbuffer"
)

const (
	instanceTableName             = "instances"
	instanceIdIndexName           = "id"
	instanceDeploymentIdIndexName = "deployment-id"
)

func init() {
	schemas = append(schemas, instanceSchema)
}

func instanceSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: instanceTableName,
		Indexes: map[string]*memdb.IndexSchema{
			instanceIdIndexName: &memdb.IndexSchema{
				Name:         instanceIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			instanceDeploymentIdIndexName: &memdb.IndexSchema{
				Name:         instanceDeploymentIdIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "DeploymentId",
					Lowercase: true,
				},
			},
		},
	}
}

type Instance struct {
	Id           string
	DeploymentId string
	LogBuffer    *logbuffer.Buffer
}

func (s *State) InstanceCreate(rec *Instance) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Create our instance
	if err := txn.Insert(instanceTableName, rec); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	// Delete all the instance exec sessions. This should be empty anyways.
	if _, err := txn.DeleteAll(instanceExecTableName, instanceExecInstanceIdIndexName, rec.Id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	txn.Commit()

	return nil
}

func (s *State) InstanceDelete(id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll(instanceTableName, instanceIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceById(id string) (*Instance, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(instanceTableName, instanceIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "instance ID not found")
	}

	return raw.(*Instance), nil
}

func (s *State) InstancesByDeployment(id string, ws memdb.WatchSet) ([]*Instance, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()
	iter, err := txn.Get(instanceTableName, instanceDeploymentIdIndexName, id)
	if err != nil {
		return nil, err
	}
	ws.Add(iter.WatchCh())

	var result []*Instance
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*Instance))
	}

	return result, nil
}
