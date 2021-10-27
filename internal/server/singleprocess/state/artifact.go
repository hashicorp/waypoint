package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/hashicorp/waypoint/internal/serverstate"
)

var artifactOp = &appOperation{
	Struct: (*pb.PushedArtifact)(nil),
	Bucket: []byte("artifact"),
}

func init() {
	artifactOp.register()
}

// ArtifactPut inserts or updates a artifact record.
func (s *State) ArtifactPut(update bool, b *pb.PushedArtifact) error {
	return artifactOp.Put(s, update, b)
}

// ArtifactGet gets a artifact by ref.
func (s *State) ArtifactGet(ref *pb.Ref_Operation) (*pb.PushedArtifact, error) {
	result, err := artifactOp.Get(s, ref)
	if err != nil {
		return nil, err
	}

	return result.(*pb.PushedArtifact), nil
}

func (s *State) ArtifactList(
	ref *pb.Ref_Application,
	opts ...serverstate.ListOperationOption,
) ([]*pb.PushedArtifact, error) {
	raw, err := artifactOp.List(s, serverstate.BuildListOperationOptions(ref, opts...))
	if err != nil {
		return nil, err
	}

	result := make([]*pb.PushedArtifact, len(raw))
	for i, v := range raw {
		result[i] = v.(*pb.PushedArtifact)
	}

	return result, nil
}

// ArtifactLatest gets the latest artifact that was completed successfully.
func (s *State) ArtifactLatest(
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (*pb.PushedArtifact, error) {
	result, err := artifactOp.Latest(s, ref, ws)
	if err != nil {
		return nil, err
	}

	return result.(*pb.PushedArtifact), nil
}
