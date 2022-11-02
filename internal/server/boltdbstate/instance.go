package boltdbstate

import (
	"context"

	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
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

func (s *State) InstanceCreate(ctx context.Context, rec *serverstate.Instance) error {
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

func (s *State) InstanceDelete(ctx context.Context, id string) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll(instanceTableName, instanceIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceById(ctx context.Context, id string) (*serverstate.Instance, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(instanceTableName, instanceIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "instance ID not found")
	}

	return raw.(*serverstate.Instance), nil
}

// instanceByIdWaiting waits for an instance with +id+ to connect before returning
// itself record.
func (s *State) instanceByIdWaiting(ctx context.Context, id string) (*serverstate.Instance, error) {
	// If the caller specified an instance id already, then just validate it.
	if id == "" {
		return nil, status.Errorf(codes.NotFound, "No instance id given")
	}

	for {
		// We have to start a new txn per iteration because we want to be able to observe
		// the newly created record for the instance.
		txn := s.inmem.Txn(false)

		// NOTE: we don't defer the txn.Abort() here because Abort() on a readonly txn
		// is a noop anyway AND we don't want to fill the stack of this function up with
		// defers, since this is in a loop. Defers in loops, thar be dragons.

		watchCh, raw, err := txn.FirstWatch(instanceTableName, instanceIdIndexName, id)
		if err != nil {
			return nil, err
		}

		// It's there!
		if raw != nil {
			return raw.(*serverstate.Instance), nil
		}

		// The watcher here registers itself on the bottom of a leaf node in the memdb
		// graph, which is fired when a new value is loaded into that leaf. Thus, it can
		// detect previously unknown values.
		ws := memdb.NewWatchSet()
		ws.Add(watchCh)

		// Wait for the instance to show up
		if err := ws.WatchCtx(ctx); err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			return nil, err
		}
	}
}

func (s *State) InstancesByDeployment(ctx context.Context, id string, ws memdb.WatchSet) ([]*serverstate.Instance, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()
	iter, err := txn.Get(instanceTableName, instanceDeploymentIdIndexName, id)
	if err != nil {
		return nil, err
	}
	ws.Add(iter.WatchCh())

	var result []*serverstate.Instance
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*serverstate.Instance))
	}

	return result, nil
}

func (s *State) InstancesByApp(
	ctx context.Context,
	ref *pb.Ref_Application,
	refws *pb.Ref_Workspace,
	ws memdb.WatchSet,
) ([]*serverstate.Instance, error) {
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

	var result []*serverstate.Instance
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*serverstate.Instance))
	}

	return result, nil
}
