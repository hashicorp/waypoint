package boltdbstate

import (
	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
	bolt "go.etcd.io/bbolt"
)

var deploymentOp = &appOperation{
	Struct: (*pb.Deployment)(nil),
	Bucket: []byte("deployment"),

	MaximumIndexedRecords: 10000,
}

func init() {
	deploymentOp.register()
}

// DeploymentPut inserts or updates a deployment record.
func (s *State) DeploymentPut(update bool, b *pb.Deployment) error {
	return deploymentOp.Put(s, update, b)
}

// DeploymentGet gets a deployment by ref.
func (s *State) DeploymentGet(ref *pb.Ref_Operation) (*pb.Deployment, error) {
	result, err := deploymentOp.Get(s, ref)
	if err != nil {
		return nil, err
	}

	return result.(*pb.Deployment), nil
}

func (s *State) DeploymentList(
	ref *pb.Ref_Application,
	opts ...serverstate.ListOperationOption,
) ([]*pb.Deployment, error) {
	raw, err := deploymentOp.List(s, serverstate.BuildListOperationOptions(ref, opts...))
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Deployment, len(raw))
	for i, v := range raw {
		result[i] = v.(*pb.Deployment)
	}

	return result, nil
}

// DeploymentLatest gets the latest deployment that was completed successfully.
func (s *State) DeploymentLatest(
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (*pb.Deployment, error) {
	result, err := deploymentOp.Latest(s, ref, ws)
	if result == nil || err != nil {
		return nil, err
	}

	return result.(*pb.Deployment), nil
}

// DeploymentDelete deletes the deployment from the DB
func (s *State) DeploymentDelete(
	ref *pb.Ref_Operation,
) error {
	return deploymentOp.Delete(s, nil)
}

func (s *State) deploymentDelete(dbTxn *bolt.Tx, memTxn *memdb.Txn, d *pb.Deployment) error {
	return deploymentOp.delete(dbTxn, memTxn, d)
}
