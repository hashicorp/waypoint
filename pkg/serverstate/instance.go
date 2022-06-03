package serverstate

import (
	"context"

	"github.com/hashicorp/go-memdb"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
)

// Instance represents a running deployment instance for an application.
type Instance struct {
	Id           string
	DeploymentId string
	Project      string
	Application  string
	Workspace    string
	LogBuffer    *logbuffer.Buffer
	Type         pb.Instance_Type
	DisableExec  bool
}

func (i *Instance) Proto() *pb.Instance {
	return &pb.Instance{
		Id:           i.Id,
		DeploymentId: i.DeploymentId,
		Type:         i.Type,
		Application: &pb.Ref_Application{
			Project:     i.Project,
			Application: i.Application,
		},
		Workspace: &pb.Ref_Workspace{
			Workspace: i.Workspace,
		},
	}
}

// InstanceExec represents a single exec session.
type InstanceExec struct {
	Id         int64
	InstanceId string

	Args []string
	Pty  *pb.ExecStreamRequest_PTY

	ClientEventCh     <-chan *pb.ExecStreamRequest
	EntrypointEventCh chan<- *pb.EntrypointExecRequest
	Connected         uint32

	// This is the context that the client side is running inside.
	// It is used by the entrypoint side to detect if the client is still
	// around or not.
	Context context.Context
}

type InstanceExecHandler interface {
	InstanceExecCreateByTargetedInstance(id string, exec *InstanceExec) error
	InstanceExecCreateByDeployment(did string, exec *InstanceExec) error
	InstanceExecCreateForVirtualInstance(ctx context.Context, id string, exec *InstanceExec) error
	InstanceExecDelete(id int64) error
	InstanceExecById(id int64) (*InstanceExec, error)
	InstanceExecConnect(ctx context.Context, id int64) (*InstanceExec, error)
	InstanceExecListByInstanceId(id string, ws memdb.WatchSet) ([]*InstanceExec, error)
	InstanceExecWaitConnected(ctx context.Context, exec *InstanceExec) error
}
