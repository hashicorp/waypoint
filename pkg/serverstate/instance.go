package serverstate

import (
	"context"

	"github.com/hashicorp/waypoint/internal/server/logbuffer"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
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
