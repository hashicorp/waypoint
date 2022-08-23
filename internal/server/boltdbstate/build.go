package boltdbstate

import (
	"errors"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverstate"
	bolt "go.etcd.io/bbolt"
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

// BuildGet gets a build by ref.
func (s *State) BuildGet(ref *pb.Ref_Operation) (*pb.Build, error) {
	result, err := buildOp.Get(s, ref)
	if err != nil {
		return nil, err
	}

	return result.(*pb.Build), nil
}

func (s *State) BuildList(
	ref *pb.Ref_Application,
	opts ...serverstate.ListOperationOption,
) ([]*pb.Build, error) {
	raw, err := buildOp.List(s, serverstate.BuildListOperationOptions(ref, opts...))
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
func (s *State) BuildLatest(
	ref *pb.Ref_Application,
	ws *pb.Ref_Workspace,
) (*pb.Build, error) {
	result, err := buildOp.Latest(s, ref, ws)
	if result == nil || err != nil {
		return nil, err
	}

	return result.(*pb.Build), nil
}

// BuildDelete deletes the build from the DB
func (s *State) BuildDelete(
	ref *pb.Ref_Operation,
) error {
	return buildOp.Delete(s, ref)
}

func (s *State) buildDelete(
	dbTxn *bolt.Tx,
	ref *pb.Ref_Operation,
) error {
	id, ok := ref.Target.(*pb.Ref_Operation_Id)
	if !ok {
		return errors.New("invalid type for target to delete app operation")
	}
	return buildOp.delete(dbTxn, []byte(id.Id))
}
