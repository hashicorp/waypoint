package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var deploymentOp = &appOperation{
	Struct: (*pb.Deployment)(nil),
	Bucket: []byte("deployment"),
}

func init() {
	deploymentOp.register()
}

// DeploymentPut inserts or updates a deployment record.
func (s *State) DeploymentPut(update bool, b *pb.Deployment) error {
	return deploymentOp.Put(s, update, b)
}

// DeploymentGet gets a deployment by ID.
func (s *State) DeploymentGet(id string) (*pb.Deployment, error) {
	result, err := deploymentOp.Get(s, id)
	if err != nil {
		return nil, err
	}

	return result.(*pb.Deployment), nil
}

func (s *State) DeploymentList(ref *pb.Ref_Application) ([]*pb.Deployment, error) {
	raw, err := deploymentOp.List(s, ref)
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
func (s *State) DeploymentLatest(ref *pb.Ref_Application) (*pb.Deployment, error) {
	result, err := deploymentOp.Latest(s, ref)
	if result == nil || err != nil {
		return nil, err
	}

	return result.(*pb.Deployment), nil
}
