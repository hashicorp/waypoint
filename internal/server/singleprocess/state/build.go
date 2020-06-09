package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var buildOp = &appOperation{
	Struct: (*pb.Build)(nil),
	Bucket: []byte("build"),
}

func init() {
	buildOp.register()
}

// BuildPut inserts or updates a build record.
func (s *State) BuildPut(update bool, b *pb.Build) error {
	return buildOp.Put(s, update, b)
}

// BuildGet gets a build by ID.
func (s *State) BuildGet(id string) (*pb.Build, error) {
	result, err := buildOp.Get(s, id)
	if err != nil {
		return nil, err
	}

	return result.(*pb.Build), nil
}

func (s *State) BuildList(ref *pb.Ref_Application) ([]*pb.Build, error) {
	raw, err := buildOp.List(s, &listOperationsOptions{
		Application: ref,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*pb.Build, len(raw))
	for i, v := range raw {
		result[i] = v.(*pb.Build)
	}

	return result, nil
}

// BuildLatest gets the latest build that was completed successfully.
func (s *State) BuildLatest(ref *pb.Ref_Application) (*pb.Build, error) {
	result, err := buildOp.Latest(s, ref)
	if result == nil || err != nil {
		return nil, err
	}

	return result.(*pb.Build), nil
}
