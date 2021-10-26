package serverstate

import (
	"context"
	"time"

	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Interface is the primary interface implemented by an implementation.
//
// This is an internal interface because (1) we don't expect or support
// any external implementations and (2) we can absolutely change this interface
// anytime we want or find it convenient, but we have to make sure so simultaneously
// modify all our implementations.
type Interface interface {
	HMACKeyEmpty() bool
	HMACKeyCreateIfNotExist(id string, size int) (*pb.HMACKey, error)
	HMACKeyGet(id string) (*pb.HMACKey, error)

	ServerConfigSet(*pb.ServerConfig) error
	ServerConfigGet() (*pb.ServerConfig, error)

	UserPut(*pb.User) error
	UserGet(*pb.Ref_User) (*pb.User, error)
	UserDelete(*pb.Ref_User) error
	UserList() ([]*pb.User, error)
	UserEmpty() (bool, error)
	UserGetOIDC(iss, sub string) (*pb.User, error)
	UserGetEmail(string) (*pb.User, error)

	//---------------------------------------------------------------
	// Server Settings

	AuthMethodPut(*pb.AuthMethod) error
	AuthMethodGet(*pb.Ref_AuthMethod) (*pb.AuthMethod, error)
	AuthMethodDelete(*pb.Ref_AuthMethod) error
	AuthMethodList() ([]*pb.AuthMethod, error)

	//---------------------------------------------------------------
	// Config (App, Runner, etc.)

	ConfigSet(...*pb.ConfigVar) error
	ConfigGet(*pb.ConfigGetRequest) ([]*pb.ConfigVar, error)
	ConfigGetWatch(*pb.ConfigGetRequest, memdb.WatchSet) ([]*pb.ConfigVar, error)

	ConfigSourceSet(...*pb.ConfigSource) error
	ConfigSourceGet(*pb.GetConfigSourceRequest) ([]*pb.ConfigSource, error)
	ConfigSourceGetWatch(*pb.GetConfigSourceRequest, memdb.WatchSet) ([]*pb.ConfigSource, error)

	//---------------------------------------------------------------
	// Instances

	InstanceCreate(*Instance) error
	InstanceDelete(string) error
	InstanceById(string) (*Instance, error)
	InstanceByIdWaiting(context.Context, string) (*Instance, error)
	InstancesByApp(*pb.Ref_Application, *pb.Ref_Workspace, memdb.WatchSet) ([]*Instance, error)

	//---------------------------------------------------------------
	// Projects, Apps, Workspaces

	ProjectPut(*pb.Project) error
	ProjectGet(*pb.Ref_Project) (*pb.Project, error)
	ProjectDelete(*pb.Ref_Project) error
	ProjectUpdateDataRef(*pb.Ref_Project, *pb.Ref_Workspace, *pb.Job_DataSource_Ref) error
	ProjectList() ([]*pb.Ref_Project, error)
	ProjectListWorkspaces(*pb.Ref_Project) ([]*pb.Workspace_Project, error)
	ProjectPollPeek(memdb.WatchSet) (*pb.Project, time.Time, error)
	ProjectPollComplete(*pb.Project, time.Time) error

	AppPut(*pb.Application) (*pb.Application, error)
	AppDelete(*pb.Ref_Application) error
	AppGet(*pb.Ref_Application) (*pb.Application, error)
	ApplicationPollPeek(memdb.WatchSet) (*pb.Project, time.Time, error)
	ApplicationPollComplete(*pb.Project, time.Time) error
	GetFileChangeSignal(*pb.Ref_Application) (string, error)

	//---------------------------------------------------------------
	// Operations

	ArtifactPut(bool, *pb.PushedArtifact) error
	ArtifactGet(*pb.Ref_Operation) (*pb.PushedArtifact, error)
	ArtifactLatest(*pb.Ref_Application, *pb.Ref_Workspace) (*pb.PushedArtifact, error)
	ArtifactList(*pb.Ref_Application, ...ListOperationOption) ([]*pb.PushedArtifact, error)

	BuildPut(bool, *pb.Build) error
	BuildGet(*pb.Ref_Operation) (*pb.Build, error)
	BuildLatest(*pb.Ref_Application, *pb.Ref_Workspace) (*pb.Build, error)
	BuildList(*pb.Ref_Application, ...ListOperationOption) ([]*pb.Build, error)

	DeploymentPut(bool, *pb.Deployment) error
	DeploymentGet(*pb.Ref_Operation) (*pb.Deployment, error)
	DeploymentLatest(*pb.Ref_Application, *pb.Ref_Workspace) (*pb.Deployment, error)
	DeploymentList(*pb.Ref_Application, ...ListOperationOption) ([]*pb.Deployment, error)

	ReleasePut(bool, *pb.Release) error
	ReleaseGet(*pb.Ref_Operation) (*pb.Release, error)
	ReleaseLatest(*pb.Ref_Application, *pb.Ref_Workspace) (*pb.Release, error)
	ReleaseList(*pb.Ref_Application, ...ListOperationOption) ([]*pb.Release, error)

	StatusReportPut(bool, *pb.StatusReport) error
	StatusReportGet(*pb.Ref_Operation) (*pb.StatusReport, error)
	StatusReportLatest(
		*pb.Ref_Application,
		*pb.Ref_Workspace,
		func(*pb.StatusReport) (bool, error),
	) (*pb.StatusReport, error)
	StatusReportList(*pb.Ref_Application, ...ListOperationOption) ([]*pb.StatusReport, error)
}
