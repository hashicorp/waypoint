package state

import (
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/server/logbuffer"
)

const (
	instanceTableName             = "instances"
	instanceIdIndexName           = "id"
	instanceDeploymentIdIndexName = "deployment-id"
	instanceAppIndexName          = "app"
	instanceAppWorkspaceIndexName = "app-ws"
)

func init() {
	schemas = append(schemas, instanceSchema)
}

func instanceSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: instanceTableName,
		Indexes: map[string]*memdb.IndexSchema{
			instanceIdIndexName: {
				Name:         instanceIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			instanceDeploymentIdIndexName: {
				Name:         instanceDeploymentIdIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "DeploymentId",
					Lowercase: true,
				},
			},

			instanceAppIndexName: {
				Name:         instanceAppIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Application",
							Lowercase: true,
						},
					},
				},
			},

			instanceAppWorkspaceIndexName: {
				Name:         instanceAppWorkspaceIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Application",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Workspace",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

type Instance struct {
	Id           string
	DeploymentId string
	Project      string
	Application  string
	Workspace    string
	LogBuffer    *logbuffer.Buffer
}

func (i *Instance) Proto() *pb.Instance {
	return &pb.Instance{
		Id:           i.Id,
		DeploymentId: i.DeploymentId,
		Application: &pb.Ref_Application{
			Project:     i.Project,
			Application: i.Application,
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: i.Workspace,
		},
	}
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

func (s *State) InstancesByApp(
	ref *pb.Ref_Application,
	refws *pb.Ref_Workspace,
	ws memdb.WatchSet,
) ([]*Instance, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()

	var iter memdb.ResultIterator
	var err error
	if refws == nil {
		iter, err = txn.Get(instanceTableName, instanceAppIndexName, ref.Project, ref.Application)
	} else {
		iter, err = txn.Get(instanceTableName, instanceAppWorkspaceIndexName,
			ref.Project, ref.Application, refws.Workspace)
	}
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
