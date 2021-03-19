package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
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
	opts ...ListOperationOption,
) ([]*pb.Deployment, error) {
	raw, err := deploymentOp.List(s, buildListOperationsOptions(ref, opts...))
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
