// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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

// InstanceExecHandler is an optional interface that the state interface can implement.
// When it does, the functionality associated with `waypoint exec` will be available.
type InstanceExecHandler interface {
	// InstanceExecCreateByTargetedInstance registers an exec request for a specific instance,
	// identified by it's database id.
	InstanceExecCreateByTargetedInstance(ctx context.Context, id string, exec *InstanceExec) error

	// InstanceExecCreateByDeployment looks up the instances running the given deployment,
	// picks an instance, and assigns the exec to that instance.
	InstanceExecCreateByDeployment(ctx context.Context, did string, exec *InstanceExec) error

	// InstanceExecCreateForVirtualInstance is used for plugins that require instances to be
	// created just to handle exec requests. The given instance id doesn't have to exist yet,
	// the code will wait for the instance to come online, then assign the exec request to it.
	InstanceExecCreateForVirtualInstance(ctx context.Context, id string, exec *InstanceExec) error

	// InstanceExecDelete deletes an exec request, identified by it's numeric id (InstanceExec.Id)
	InstanceExecDelete(ctx context.Context, id int64) error

	// InstanceExecById retrieves an exec request, identified by it's numeric id (InstanceExec.Id)
	InstanceExecById(ctx context.Context, id int64) (*InstanceExec, error)

	// InstanceExecListByInstanceId returns any exec requests for the given instance id.
	InstanceExecListByInstanceId(ctx context.Context, id string, ws memdb.WatchSet) ([]*InstanceExec, error)

	// InstanceExecById retrieves an exec request, identified by it's numeric id (InstanceExec.Id).
	// The implementer also can use this opertunity to do any sychronization of state, such as
	// connecting to external systems.
	InstanceExecConnect(ctx context.Context, id int64) (*InstanceExec, error)

	// InstanceExecWaitConnected is run to allow the implementation to sync up with a call to
	// InstanceExecConnect run elsewhere.
	InstanceExecWaitConnected(ctx context.Context, exec *InstanceExec) error
}
