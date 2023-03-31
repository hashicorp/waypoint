// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

	HMACKeyEmpty(context.Context) bool
	HMACKeyCreateIfNotExist(ctx context.Context, id string, size int) (*pb.HMACKey, error)
	HMACKeyGet(ctx context.Context, id string) (*pb.HMACKey, error)
	TokenSignature(ctx context.Context, tokenBody []byte, keyId string) (signature []byte, err error)
	TokenSignatureVerify(ctx context.Context, tokenBody []byte, signature []byte, keyId string) (isValid bool, err error)

	ServerConfigSet(context.Context, *pb.ServerConfig) error
	ServerConfigGet(context.Context) (*pb.ServerConfig, error)

	UserPut(context.Context, *pb.User) error
	UserGet(context.Context, *pb.Ref_User) (*pb.User, error)
	UserDelete(context.Context, *pb.Ref_User) error
	UserList(context.Context) ([]*pb.User, error)
	UserEmpty(context.Context) (bool, error)
	UserGetOIDC(ctx context.Context, iss, sub string) (*pb.User, error)
	UserGetEmail(context.Context, string) (*pb.User, error)

	//---------------------------------------------------------------
	// Server Settings

	AuthMethodPut(context.Context, *pb.AuthMethod) error
	AuthMethodGet(context.Context, *pb.Ref_AuthMethod) (*pb.AuthMethod, error)
	AuthMethodDelete(context.Context, *pb.Ref_AuthMethod) error
	AuthMethodList(context.Context) ([]*pb.AuthMethod, error)

	RunnerCreate(context.Context, *pb.Runner) error
	RunnerDelete(context.Context, string) error
	RunnerOffline(context.Context, string) error
	RunnerAdopt(context.Context, string, bool) error
	RunnerReject(context.Context, string) error
	RunnerById(context.Context, string, memdb.WatchSet) (*pb.Runner, error)
	RunnerList(context.Context) ([]*pb.Runner, error)

	OnDemandRunnerConfigPut(context.Context, *pb.OnDemandRunnerConfig) (*pb.OnDemandRunnerConfig, error)
	OnDemandRunnerConfigGet(context.Context, *pb.Ref_OnDemandRunnerConfig) (*pb.OnDemandRunnerConfig, error)
	OnDemandRunnerConfigDelete(context.Context, *pb.Ref_OnDemandRunnerConfig) error
	OnDemandRunnerConfigList(context.Context) ([]*pb.OnDemandRunnerConfig, error)
	OnDemandRunnerConfigDefault(context.Context) ([]*pb.Ref_OnDemandRunnerConfig, error)

	ServerURLTokenSet(context.Context, string) error
	ServerURLTokenGet(context.Context) (string, error)

	ServerIdSet(ctx context.Context, id string) error
	ServerIdGet(context.Context) (string, error)

	CreateSnapshot(context.Context, io.Writer) error
	StageRestoreSnapshot(context.Context, io.Reader) error

	//---------------------------------------------------------------
	// Config (App, Runner, etc.)

	ConfigSet(context.Context, ...*pb.ConfigVar) error
	ConfigGet(context.Context, *pb.ConfigGetRequest) ([]*pb.ConfigVar, error)
	ConfigGetWatch(context.Context, *pb.ConfigGetRequest, memdb.WatchSet) ([]*pb.ConfigVar, error)

	ConfigSourceSet(context.Context, ...*pb.ConfigSource) error
	ConfigSourceGet(context.Context, *pb.GetConfigSourceRequest) ([]*pb.ConfigSource, error)
	ConfigSourceGetWatch(context.Context, *pb.GetConfigSourceRequest, memdb.WatchSet) ([]*pb.ConfigSource, error)

	//---------------------------------------------------------------
	// Instances

	InstanceCreate(context.Context, *Instance) error
	InstanceDelete(context.Context, string) error
	InstanceById(context.Context, string) (*Instance, error)
	InstancesByApp(context.Context, *pb.Ref_Application, *pb.Ref_Workspace, memdb.WatchSet) ([]*Instance, error)
	InstancesByDeployment(context.Context, string, memdb.WatchSet) ([]*Instance, error)

	//---------------------------------------------------------------
	// Projects, Apps, Workspaces

	WorkspaceList(context.Context) ([]*pb.Workspace, error)
	WorkspaceListByProject(context.Context, *pb.Ref_Project) ([]*pb.Workspace, error)
	WorkspaceListByApp(context.Context, *pb.Ref_Application) ([]*pb.Workspace, error)
	WorkspaceGet(context.Context, string) (*pb.Workspace, error)
	WorkspacePut(context.Context, *pb.Workspace) error
	WorkspaceDelete(context.Context, string) error

	ProjectPut(context.Context, *pb.Project) error
	ProjectGet(context.Context, *pb.Ref_Project) (*pb.Project, error)
	ProjectDelete(context.Context, *pb.Ref_Project) error
	ProjectUpdateDataRef(context.Context, *pb.Ref_Project, *pb.Ref_Workspace, *pb.Job_DataSource_Ref) error
	ProjectCount(context.Context) (uint64, error)
	ProjectList(context.Context, *pb.PaginationRequest) ([]*pb.Ref_Project, *pb.PaginationResponse, error)
	ProjectListBundles(context.Context, *pb.PaginationRequest) ([]*pb.UI_ProjectBundle, *pb.PaginationResponse, error)
	ProjectListWorkspaces(context.Context, *pb.Ref_Project) ([]*pb.Workspace_Project, error)
	ProjectPollPeek(context.Context, memdb.WatchSet) (*pb.Project, time.Time, error)
	ProjectPollComplete(context.Context, *pb.Project, time.Time) error

	AppPut(context.Context, *pb.Application) (*pb.Application, error)
	AppDelete(context.Context, *pb.Ref_Application) error
	AppGet(context.Context, *pb.Ref_Application) (*pb.Application, error)
	ApplicationPollPeek(context.Context, memdb.WatchSet) (*pb.Project, time.Time, error)
	ApplicationPollComplete(context.Context, *pb.Project, time.Time) error
	GetFileChangeSignal(context.Context, *pb.Ref_Application) (string, error)

	//---------------------------------------------------------------
	// Operations

	ArtifactPut(context.Context, bool, *pb.PushedArtifact) error
	ArtifactGet(context.Context, *pb.Ref_Operation) (*pb.PushedArtifact, error)
	ArtifactLatest(context.Context, *pb.Ref_Application, *pb.Ref_Workspace) (*pb.PushedArtifact, error)
	ArtifactList(context.Context, *pb.Ref_Application, ...ListOperationOption) ([]*pb.PushedArtifact, error)

	BuildPut(context.Context, bool, *pb.Build) error
	BuildGet(context.Context, *pb.Ref_Operation) (*pb.Build, error)
	BuildLatest(context.Context, *pb.Ref_Application, *pb.Ref_Workspace) (*pb.Build, error)
	BuildList(context.Context, *pb.Ref_Application, ...ListOperationOption) ([]*pb.Build, error)

	DeploymentPut(context.Context, bool, *pb.Deployment) error
	DeploymentGet(context.Context, *pb.Ref_Operation) (*pb.Deployment, error)
	DeploymentLatest(context.Context, *pb.Ref_Application, *pb.Ref_Workspace) (*pb.Deployment, error)
	DeploymentList(context.Context, *pb.Ref_Application, ...ListOperationOption) ([]*pb.Deployment, error)

	ReleasePut(context.Context, bool, *pb.Release) error
	ReleaseGet(context.Context, *pb.Ref_Operation) (*pb.Release, error)
	ReleaseLatest(context.Context, *pb.Ref_Application, *pb.Ref_Workspace) (*pb.Release, error)
	ReleaseList(context.Context, *pb.Ref_Application, ...ListOperationOption) ([]*pb.Release, error)

	StatusReportPut(context.Context, bool, *pb.StatusReport) error
	StatusReportGet(context.Context, *pb.Ref_Operation) (*pb.StatusReport, error)
	StatusReportLatest(
		context.Context,
		*pb.Ref_Application,
		*pb.Ref_Workspace,
		func(*pb.StatusReport) (bool, error),
	) (*pb.StatusReport, error)
	StatusReportList(context.Context, *pb.Ref_Application, ...ListOperationOption) ([]*pb.StatusReport, error)

	//---------------------------------------------------------------
	// Trigger

	TriggerPut(context.Context, *pb.Trigger) error
	TriggerGet(context.Context, *pb.Ref_Trigger) (*pb.Trigger, error)
	TriggerDelete(context.Context, *pb.Ref_Trigger) error
	TriggerList(context.Context, *pb.Ref_Workspace, *pb.Ref_Project, *pb.Ref_Application, []string) ([]*pb.Trigger, error)

	//---------------------------------------------------------------
	// Job System

	JobCreate(context.Context, ...*pb.Job) error
	JobProjectScopedRequest(context.Context, *pb.Ref_Project, *pb.Job) ([]*pb.QueueJobRequest, error)
	JobList(context.Context, *pb.ListJobsRequest) ([]*pb.Job, *pb.PaginationResponse, error)
	JobLatestInit(context.Context, *pb.Ref_Project) (*pb.Job, error)
	JobById(context.Context, string, memdb.WatchSet) (*Job, error)
	JobPeekForRunner(context.Context, *pb.Runner) (*Job, error)
	JobAssignForRunner(context.Context, *pb.Runner) (*Job, error)
	JobAck(context.Context, string, bool) (*Job, error)
	JobUpdateRef(context.Context, string, *pb.Job_DataSource_Ref) error
	JobUpdateExpiry(context.Context, string, *timestamppb.Timestamp) error
	JobUpdate(context.Context, string, func(*pb.Job) error) error
	JobComplete(context.Context, string, *pb.Job_Result, error) error
	JobCancel(context.Context, string, bool) error
	JobHeartbeat(context.Context, string) error
	JobExpire(context.Context, string) error
	JobIsAssignable(context.Context, *pb.Job) (bool, error)

	//---------------------------------------------------------------
	// Task Tracking

	TaskPut(context.Context, *pb.Task) error
	TaskGet(context.Context, *pb.Ref_Task) (*pb.Task, error)
	TaskDelete(context.Context, *pb.Ref_Task) error
	TaskCancel(context.Context, *pb.Ref_Task) error
	TaskList(context.Context, *pb.ListTaskRequest) ([]*pb.Task, error)
	JobsByTaskRef(context.Context, *pb.Task) (*pb.Job, *pb.Job, *pb.Job, *pb.Job, error)

	//---------------------------------------------------------------
	// Pipelines

	PipelinePut(context.Context, *pb.Pipeline) error
	PipelineGet(context.Context, *pb.Ref_Pipeline) (*pb.Pipeline, error)
	PipelineDelete(context.Context, *pb.Ref_Pipeline) error
	PipelineList(context.Context, *pb.Ref_Project) ([]*pb.Pipeline, error)

	PipelineRunPut(context.Context, *pb.PipelineRun) error
	PipelineRunGet(context.Context, *pb.Ref_Pipeline, uint64) (*pb.PipelineRun, error)
	PipelineRunGetLatest(context.Context, string) (*pb.PipelineRun, error)
	PipelineRunGetById(context.Context, string) (*pb.PipelineRun, error)
	PipelineRunList(context.Context, *pb.Ref_Pipeline) ([]*pb.PipelineRun, error)

	//---------------------------------------------------------------
	// Templates
	ProjectTemplatePut(context.Context, *pb.ProjectTemplate) error
	ProjectTemplateGet(context.Context, *pb.Ref_ProjectTemplate) (*pb.ProjectTemplate, error)
	ProjectTemplateDelete(context.Context, *pb.Ref_ProjectTemplate) error
	ProjectTemplateList(context.Context, *pb.ListProjectTemplatesRequest) ([]*pb.ProjectTemplate, error)
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
