package serverstate

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-memdb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// Interface is the primary interface implemented by an implementation.
//
// Any changes to this interface will require changes to all
// implementations in all projects.
type Interface interface {
	// Close is always called when the server is shutting down or reloading
	// the state store. This should clean up any resources (file handles, etc.)
	// that the state storage is using.
	io.Closer

	HMACKeyEmpty() bool
	HMACKeyCreateIfNotExist(id string, size int) (*pb.HMACKey, error)
	HMACKeyGet(id string) (*pb.HMACKey, error)
	TokenSignature(tokenBody []byte, keyId string) (signature []byte, err error)
	TokenSignatureVerify(tokenBody []byte, signature []byte, keyId string) (isValid bool, err error)

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

	RunnerCreate(*pb.Runner) error
	RunnerDelete(string) error
	RunnerOffline(string) error
	RunnerAdopt(string, bool) error
	RunnerReject(string) error
	RunnerById(string, memdb.WatchSet) (*pb.Runner, error)
	RunnerList() ([]*pb.Runner, error)

	OnDemandRunnerConfigPut(*pb.OnDemandRunnerConfig) (*pb.OnDemandRunnerConfig, error)
	OnDemandRunnerConfigGet(*pb.Ref_OnDemandRunnerConfig) (*pb.OnDemandRunnerConfig, error)
	OnDemandRunnerConfigDelete(*pb.Ref_OnDemandRunnerConfig) error
	OnDemandRunnerConfigList() ([]*pb.OnDemandRunnerConfig, error)
	OnDemandRunnerConfigDefault() ([]*pb.Ref_OnDemandRunnerConfig, error)

	ServerURLTokenSet(string) error
	ServerURLTokenGet() (string, error)

	ServerIdSet(id string) error
	ServerIdGet() (string, error)

	CreateSnapshot(io.Writer) error
	StageRestoreSnapshot(io.Reader) error

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
	InstancesByApp(*pb.Ref_Application, *pb.Ref_Workspace, memdb.WatchSet) ([]*Instance, error)
	InstancesByDeployment(string, memdb.WatchSet) ([]*Instance, error)

	//---------------------------------------------------------------
	// Projects, Apps, Workspaces

	WorkspaceList() ([]*pb.Workspace, error)
	WorkspaceListByProject(*pb.Ref_Project) ([]*pb.Workspace, error)
	WorkspaceListByApp(*pb.Ref_Application) ([]*pb.Workspace, error)
	WorkspaceGet(string) (*pb.Workspace, error)
	WorkspacePut(*pb.Workspace) error
	WorkspaceDelete(string) error

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

	//---------------------------------------------------------------
	// Trigger

	TriggerPut(*pb.Trigger) error
	TriggerGet(*pb.Ref_Trigger) (*pb.Trigger, error)
	TriggerDelete(*pb.Ref_Trigger) error
	TriggerList(*pb.Ref_Workspace, *pb.Ref_Project, *pb.Ref_Application, []string) ([]*pb.Trigger, error)

	//---------------------------------------------------------------
	// Job System

	JobCreate(...*pb.Job) error
	JobProjectScopedRequest(*pb.Ref_Project, *pb.Job) ([]*pb.QueueJobRequest, error)
	JobList(*pb.ListJobsRequest) ([]*pb.Job, error)
	JobById(string, memdb.WatchSet) (*Job, error)
	JobPeekForRunner(context.Context, *pb.Runner) (*Job, error)
	JobAssignForRunner(context.Context, *pb.Runner) (*Job, error)
	JobAck(string, bool) (*Job, error)
	JobUpdateRef(string, *pb.Job_DataSource_Ref) error
	JobUpdateExpiry(string, *timestamppb.Timestamp) error
	JobUpdate(string, func(*pb.Job) error) error
	JobComplete(string, *pb.Job_Result, error) error
	JobCancel(string, bool) error
	JobHeartbeat(string) error
	JobExpire(string) error
	JobIsAssignable(context.Context, *pb.Job) (bool, error)

	//---------------------------------------------------------------
	// Task Tracking

	TaskPut(*pb.Task) error
	TaskGet(*pb.Ref_Task) (*pb.Task, error)
	TaskDelete(*pb.Ref_Task) error
	TaskCancel(*pb.Ref_Task) error
	TaskList(*pb.ListTaskRequest) ([]*pb.Task, error)
	JobsByTaskRef(*pb.Task) (*pb.Job, *pb.Job, *pb.Job, *pb.Job, error)

	//---------------------------------------------------------------
	// Pipelines

	PipelinePut(*pb.Pipeline) error
	PipelineGet(*pb.Ref_Pipeline) (*pb.Pipeline, error)
	PipelineDelete(*pb.Ref_Pipeline) error
	PipelineList(*pb.Ref_Project) ([]*pb.Pipeline, error)

	PipelineRunPut(*pb.PipelineRun) error
	PipelineRunGet(*pb.Ref_Pipeline, uint64) (*pb.PipelineRun, error)
	PipelineRunGetLatest(string) (*pb.PipelineRun, error)
	PipelineRunGetById(string) (*pb.PipelineRun, error)
	PipelineRunList(*pb.Ref_Pipeline) ([]*pb.PipelineRun, error)
}

// Pruner is implemented by state storage implementations that require
// a periodic prune. The implementation can't control when this is called,
// but it will be called roughly every hour or shorter.
//
// During pruning, other operations are still allowed to come in. It is up
// to the state implementation to handle safe concurrency.
type Pruner interface {
	Interface

	Prune() error
}
