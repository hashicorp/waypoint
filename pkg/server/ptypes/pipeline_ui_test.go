// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func TestUI_PipelineRunTreeFromJobs(t *testing.T) {
	appRef := &pb.Ref_Application{
		Project:     "test-project",
		Application: "test-app",
	}

	cases := map[string]struct {
		Jobs          []*pb.Job
		StatusReports []*pb.StatusReport
		Tree          *pb.UI_PipelineRunTreeNode
	}{
		"one queued exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State: pb.Job_QUEUED,
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State: pb.UI_PipelineRunTreeNode_QUEUED,
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one running exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:   pb.Job_RUNNING,
					AckTime: quickTimestamp("2023-01-01T13:00:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:     pb.UI_PipelineRunTreeNode_RUNNING,
				StartTime: quickTimestamp("2023-01-01T13:00:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one successful exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:        pb.Job_SUCCESS,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_SUCCESS,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one errored exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:        pb.Job_ERROR,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_ERROR,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one cancelled exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:        pb.Job_ERROR,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CancelTime:   quickTimestamp("2023-01-01T13:08:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_CANCELLED,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one queued exec step with application": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State:       pb.Job_QUEUED,
					Application: appRef,
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State: pb.UI_PipelineRunTreeNode_QUEUED,
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Application: appRef,
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one queued exec step with workspace": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					State: pb.Job_QUEUED,
					Workspace: &pb.Ref_Workspace{
						Workspace: "test",
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State: pb.UI_PipelineRunTreeNode_QUEUED,
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Workspace: &pb.Ref_Workspace{
					Workspace: "test",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one successful build step with result": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-build-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "build",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Build_{
									Build: &pb.Pipeline_Step_Build{
										DisablePush: true,
									},
								},
							},
						},
					},
					State:        pb.Job_SUCCESS,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
					Result: &pb.Job_Result{
						Build: &pb.Job_BuildResult{
							Build: &pb.Build{
								Id: "build-from-build-step",
							},
						},
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "build",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "build",
					Kind: &pb.Pipeline_Step_Build_{
						Build: &pb.Pipeline_Step_Build{
							DisablePush: true,
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_SUCCESS,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-build-step",
				},
				Result: &pb.Job_Result{
					Build: &pb.Job_BuildResult{
						Build: &pb.Build{
							Id: "build-from-build-step",
						},
					},
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one successful deploy step with result": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-deploy-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "deploy",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Deploy_{
									Deploy: &pb.Pipeline_Step_Deploy{
										Release: false,
									},
								},
							},
						},
					},
					State:        pb.Job_SUCCESS,
					AckTime:      quickTimestamp("2023-01-01T13:00:00Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
					Result: &pb.Job_Result{
						Deploy: &pb.Job_DeployResult{
							Deployment: &pb.Deployment{
								Id: "deployment-from-deploy-step",
							},
						},
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "deploy",
						RunSequence:  1,
					},
				},
			},
			StatusReports: []*pb.StatusReport{
				{
					TargetId: &pb.StatusReport_DeploymentId{DeploymentId: "deployment-from-deploy-step"},
					Id:       "status-report-for-deployment",
					Health:   &pb.StatusReport_Health{HealthStatus: "READY"},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "deploy",
					Kind: &pb.Pipeline_Step_Deploy_{
						Deploy: &pb.Pipeline_Step_Deploy{
							Release: false,
						},
					},
				},
				State:        pb.UI_PipelineRunTreeNode_SUCCESS,
				StartTime:    quickTimestamp("2023-01-01T13:00:00Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:10:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-deploy-step",
				},
				Result: &pb.Job_Result{
					Deploy: &pb.Job_DeployResult{
						Deployment: &pb.Deployment{
							Id: "deployment-from-deploy-step",
						},
					},
				},
				LatestStatusReport: &pb.StatusReport{
					TargetId: &pb.StatusReport_DeploymentId{DeploymentId: "deployment-from-deploy-step"},
					Id:       "status-report-for-deployment",
					Health:   &pb.StatusReport_Health{HealthStatus: "READY"},
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{},
				},
			},
		},
		"one running exec step and one queued exec step": {
			Jobs: []*pb.Job{
				{
					Id: "job-for-hello-step",
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "hello",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hello"},
									},
								},
							},
						},
					},
					AckTime: quickTimestamp("2023-01-01T13:00:00Z"),
					State:   pb.Job_RUNNING,
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "hello",
						RunSequence:  1,
					},
				},
				{
					Id:        "job-for-bye-step",
					DependsOn: []string{"job-for-hello-step"},
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "bye",
								DependsOn: []string{"hello"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"bye"},
									},
								},
							},
						},
					},
					State: pb.Job_QUEUED,
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "pipeline-id",
						PipelineName: "pipeline-name",
						Step:         "bye",
						RunSequence:  1,
					}},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name: "hello",
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"hello"},
						},
					},
				},
				State:     pb.UI_PipelineRunTreeNode_RUNNING,
				StartTime: quickTimestamp("2023-01-01T13:00:00Z"),
				Job: &pb.Ref_Job{
					Id: "job-for-hello-step",
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode: pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{
						{
							Step: &pb.Pipeline_Step{
								Name: "bye",
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"bye"},
									},
								},
								DependsOn: []string{"hello"},
							},
							State: pb.UI_PipelineRunTreeNode_QUEUED,
							Job:   &pb.Ref_Job{Id: "job-for-bye-step"},
							Children: &pb.UI_PipelineRunTreeNode_Children{
								Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
								Nodes: []*pb.UI_PipelineRunTreeNode{},
							},
						},
					},
				},
			},
		},
		"referenced pipeline invocation": {
			Jobs: []*pb.Job{
				{
					Id:          "job-for-prep-step",
					DependsOn:   []string{"unknown-job"},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "prep",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"preparing"},
									},
								},
							},
						},
					},
					QueueTime:    quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:      quickTimestamp("2023-01-01T13:00:10Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:00:20Z"),
					State:        pb.Job_SUCCESS,
					Result: &pb.Job_Result{
						PipelineStep: &pb.Job_PipelineStepResult{
							Result: &status.Status{},
						},
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "prep",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-invoke-step",
					DependsOn: []string{
						"job-for-prep-step",
						"job-for-hi-step",
						"job-for-bye-step",
						"unknown-job",
					},
					Application: appRef,
					Operation:   &pb.Job_Noop_{},
					State:       pb.Job_QUEUED,
					QueueTime:   quickTimestamp("2023-01-01T13:00:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "invoke-referenced-pipeline",
						RunSequence:  1,
					},
				},
				{
					Id:          "job-for-hi-step",
					DependsOn:   []string{"job-for-prep-step"},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name: "hi",
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hi"},
									},
								},
							},
						},
					},
					State:        pb.Job_SUCCESS,
					QueueTime:    quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:      quickTimestamp("2023-01-01T13:00:30Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:00:40Z"),
					Result: &pb.Job_Result{
						PipelineStep: &pb.Job_PipelineStepResult{
							Result: &status.Status{},
						},
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "referenced-pipeline-id",
						PipelineName: "referenced-pipeline",
						Step:         "hi",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-bye-step",
					DependsOn: []string{
						"job-for-hi-step",
						"job-for-prep-step",
						"unknown-job",
					},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "bye",
								DependsOn: []string{"hi"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"bye"},
									},
								},
							},
						},
					},
					State:     pb.Job_RUNNING,
					QueueTime: quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:   quickTimestamp("2023-01-01T13:00:50Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "referenced-pipeline-id",
						PipelineName: "referenced-pipeline",
						Step:         "bye",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-done-step",
					DependsOn: []string{
						"job-for-hi-step",
						"job-for-bye-step",
						"job-for-invoke-step",
						"job-for-prep-step",
						"unknown-job",
					},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "done",
								DependsOn: []string{"invoke-referenced-pipeline"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"done"},
									},
								},
							},
						},
					},
					State:     pb.Job_QUEUED,
					QueueTime: quickTimestamp("2023-01-01T13:00:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "done",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name:      "prep",
					DependsOn: []string{},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"preparing"},
						},
					},
				},
				Application:  appRef,
				Job:          &pb.Ref_Job{Id: "job-for-prep-step"},
				StartTime:    quickTimestamp("2023-01-01T13:00:10Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:00:20Z"),
				State:        pb.UI_PipelineRunTreeNode_SUCCESS,
				Result: &pb.Job_Result{
					PipelineStep: &pb.Job_PipelineStepResult{
						Result: &status.Status{},
					},
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode: pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{
						{
							Step: &pb.Pipeline_Step{
								Name: "invoke-referenced-pipeline",
								Kind: &pb.Pipeline_Step_Pipeline_{
									Pipeline: &pb.Pipeline_Step_Pipeline{
										Ref: &pb.Ref_Pipeline{
											Ref: &pb.Ref_Pipeline_Owner{
												Owner: &pb.Ref_PipelineOwner{
													Project: &pb.Ref_Project{
														Project: "test-project",
													},
													PipelineName: "referenced-pipeline",
												},
											},
										},
									},
								},
							},
							Application: appRef,
							Job:         &pb.Ref_Job{Id: "job-for-invoke-step"},
							State:       pb.UI_PipelineRunTreeNode_RUNNING,
							StartTime:   quickTimestamp("2023-01-01T13:00:30Z"),
							Children: &pb.UI_PipelineRunTreeNode_Children{
								Mode: pb.UI_PipelineRunTreeNode_Children_SERIAL,
								Nodes: []*pb.UI_PipelineRunTreeNode{
									{
										Step: &pb.Pipeline_Step{
											Name: "hi",
											Kind: &pb.Pipeline_Step_Exec_{
												Exec: &pb.Pipeline_Step_Exec{
													Image:   "busybox",
													Command: "echo",
													Args:    []string{"hi"},
												},
											},
										},
										Application:  appRef,
										Job:          &pb.Ref_Job{Id: "job-for-hi-step"},
										StartTime:    quickTimestamp("2023-01-01T13:00:30Z"),
										CompleteTime: quickTimestamp("2023-01-01T13:00:40Z"),
										State:        pb.UI_PipelineRunTreeNode_SUCCESS,
										Result: &pb.Job_Result{
											PipelineStep: &pb.Job_PipelineStepResult{
												Result: &status.Status{},
											},
										},
										Children: &pb.UI_PipelineRunTreeNode_Children{
											Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
											Nodes: []*pb.UI_PipelineRunTreeNode{},
										},
									},
									{
										Step: &pb.Pipeline_Step{
											Name:      "bye",
											DependsOn: []string{"hi"},
											Kind: &pb.Pipeline_Step_Exec_{
												Exec: &pb.Pipeline_Step_Exec{
													Image:   "busybox",
													Command: "echo",
													Args:    []string{"bye"},
												},
											},
										},
										Application: appRef,
										Job:         &pb.Ref_Job{Id: "job-for-bye-step"},
										StartTime:   quickTimestamp("2023-01-01T13:00:50Z"),
										State:       pb.UI_PipelineRunTreeNode_RUNNING,
										Children: &pb.UI_PipelineRunTreeNode_Children{
											Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
											Nodes: []*pb.UI_PipelineRunTreeNode{},
										},
									},
								},
							},
						},
						{
							Step: &pb.Pipeline_Step{
								Name:      "done",
								DependsOn: []string{"invoke-referenced-pipeline"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"done"},
									},
								},
							},
							Application: appRef,
							Job:         &pb.Ref_Job{Id: "job-for-done-step"},
							State:       pb.UI_PipelineRunTreeNode_QUEUED,
							Children: &pb.UI_PipelineRunTreeNode_Children{
								Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
								Nodes: []*pb.UI_PipelineRunTreeNode{},
							},
						},
					},
				},
			},
		},
		"referenced pipeline invocation AND parallel steps": {
			Jobs: []*pb.Job{
				{
					Id:          "job-for-prep-step",
					DependsOn:   []string{"unknown-job"},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "prep",
								DependsOn: []string{},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"preparing"},
									},
								},
							},
						},
					},
					QueueTime:    quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:      quickTimestamp("2023-01-01T13:00:10Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:00:20Z"),
					State:        pb.Job_SUCCESS,
					Result: &pb.Job_Result{
						PipelineStep: &pb.Job_PipelineStepResult{
							Result: &status.Status{},
						},
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "prep",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-invoke-step",
					DependsOn: []string{
						"job-for-prep-step",
						"job-for-hi-step",
						"job-for-bye-step",
						"unknown-job",
					},
					Application: appRef,
					Operation:   &pb.Job_Noop_{},
					State:       pb.Job_QUEUED,
					QueueTime:   quickTimestamp("2023-01-01T13:00:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "invoke-referenced-pipeline",
						RunSequence:  1,
					},
				},
				{
					Id:          "job-for-hi-step",
					DependsOn:   []string{"job-for-prep-step"},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name: "hi",
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"hi"},
									},
								},
							},
						},
					},
					State:        pb.Job_SUCCESS,
					QueueTime:    quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:      quickTimestamp("2023-01-01T13:00:30Z"),
					CompleteTime: quickTimestamp("2023-01-01T13:00:40Z"),
					Result: &pb.Job_Result{
						PipelineStep: &pb.Job_PipelineStepResult{
							Result: &status.Status{},
						},
					},
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "referenced-pipeline-id",
						PipelineName: "referenced-pipeline",
						Step:         "hi",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-bye-step",
					DependsOn: []string{
						"job-for-hi-step",
						"job-for-prep-step",
						"unknown-job",
					},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "bye",
								DependsOn: []string{"hi"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"bye"},
									},
								},
							},
						},
					},
					State:     pb.Job_RUNNING,
					QueueTime: quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:   quickTimestamp("2023-01-01T13:00:50Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "referenced-pipeline-id",
						PipelineName: "referenced-pipeline",
						Step:         "bye",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-other-step",
					DependsOn: []string{
						"job-for-prep-step",
						"unknown-job",
					},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "other",
								DependsOn: []string{"prep"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"other"},
									},
								},
							},
						},
					},
					State:     pb.Job_RUNNING,
					QueueTime: quickTimestamp("2023-01-01T13:00:00Z"),
					AckTime:   quickTimestamp("2023-01-01T13:00:50Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "other",
						RunSequence:  1,
					},
				},
				{
					Id: "job-for-done-step",
					DependsOn: []string{
						"job-for-hi-step",
						"job-for-bye-step",
						"job-for-invoke-step",
						"job-for-prep-step",
						"job-for-other-step",
						"unknown-job",
					},
					Application: appRef,
					Operation: &pb.Job_PipelineStep{
						PipelineStep: &pb.Job_PipelineStepOp{
							Step: &pb.Pipeline_Step{
								Name:      "done",
								DependsOn: []string{"invoke-referenced-pipeline"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"done"},
									},
								},
							},
						},
					},
					State:     pb.Job_QUEUED,
					QueueTime: quickTimestamp("2023-01-01T13:00:00Z"),
					Pipeline: &pb.Ref_PipelineStep{
						PipelineId:   "parent-pipeline-id",
						PipelineName: "parent-pipeline",
						Step:         "done",
						RunSequence:  1,
					},
				},
			},
			Tree: &pb.UI_PipelineRunTreeNode{
				Step: &pb.Pipeline_Step{
					Name:      "prep",
					DependsOn: []string{},
					Kind: &pb.Pipeline_Step_Exec_{
						Exec: &pb.Pipeline_Step_Exec{
							Image:   "busybox",
							Command: "echo",
							Args:    []string{"preparing"},
						},
					},
				},
				Application:  appRef,
				Job:          &pb.Ref_Job{Id: "job-for-prep-step"},
				StartTime:    quickTimestamp("2023-01-01T13:00:10Z"),
				CompleteTime: quickTimestamp("2023-01-01T13:00:20Z"),
				State:        pb.UI_PipelineRunTreeNode_SUCCESS,
				Result: &pb.Job_Result{
					PipelineStep: &pb.Job_PipelineStepResult{
						Result: &status.Status{},
					},
				},
				Children: &pb.UI_PipelineRunTreeNode_Children{
					Mode: pb.UI_PipelineRunTreeNode_Children_SERIAL,
					Nodes: []*pb.UI_PipelineRunTreeNode{
						{
							Job:         &pb.Ref_Job{Id: "job-for-prep-step-virtual"},
							Application: appRef,
							StartTime:   quickTimestamp("2023-01-01T13:00:30Z"),
							State:       pb.UI_PipelineRunTreeNode_RUNNING,
							Children: &pb.UI_PipelineRunTreeNode_Children{
								Mode: pb.UI_PipelineRunTreeNode_Children_PARALLEL,
								Nodes: []*pb.UI_PipelineRunTreeNode{
									{
										Step: &pb.Pipeline_Step{
											Name: "invoke-referenced-pipeline",
											Kind: &pb.Pipeline_Step_Pipeline_{
												Pipeline: &pb.Pipeline_Step_Pipeline{
													Ref: &pb.Ref_Pipeline{
														Ref: &pb.Ref_Pipeline_Owner{
															Owner: &pb.Ref_PipelineOwner{
																Project: &pb.Ref_Project{
																	Project: "test-project",
																},
																PipelineName: "referenced-pipeline",
															},
														},
													},
												},
											},
										},
										Application: appRef,
										Job:         &pb.Ref_Job{Id: "job-for-invoke-step"},
										State:       pb.UI_PipelineRunTreeNode_RUNNING,
										StartTime:   quickTimestamp("2023-01-01T13:00:30Z"),
										Children: &pb.UI_PipelineRunTreeNode_Children{
											Mode: pb.UI_PipelineRunTreeNode_Children_SERIAL,
											Nodes: []*pb.UI_PipelineRunTreeNode{
												{
													Step: &pb.Pipeline_Step{
														Name: "hi",
														Kind: &pb.Pipeline_Step_Exec_{
															Exec: &pb.Pipeline_Step_Exec{
																Image:   "busybox",
																Command: "echo",
																Args:    []string{"hi"},
															},
														},
													},
													Application:  appRef,
													Job:          &pb.Ref_Job{Id: "job-for-hi-step"},
													StartTime:    quickTimestamp("2023-01-01T13:00:30Z"),
													CompleteTime: quickTimestamp("2023-01-01T13:00:40Z"),
													State:        pb.UI_PipelineRunTreeNode_SUCCESS,
													Result: &pb.Job_Result{
														PipelineStep: &pb.Job_PipelineStepResult{
															Result: &status.Status{},
														},
													},
													Children: &pb.UI_PipelineRunTreeNode_Children{
														Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
														Nodes: []*pb.UI_PipelineRunTreeNode{},
													},
												},
												{
													Step: &pb.Pipeline_Step{
														Name:      "bye",
														DependsOn: []string{"hi"},
														Kind: &pb.Pipeline_Step_Exec_{
															Exec: &pb.Pipeline_Step_Exec{
																Image:   "busybox",
																Command: "echo",
																Args:    []string{"bye"},
															},
														},
													},
													Application: appRef,
													Job:         &pb.Ref_Job{Id: "job-for-bye-step"},
													StartTime:   quickTimestamp("2023-01-01T13:00:50Z"),
													State:       pb.UI_PipelineRunTreeNode_RUNNING,
													Children: &pb.UI_PipelineRunTreeNode_Children{
														Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
														Nodes: []*pb.UI_PipelineRunTreeNode{},
													},
												},
											},
										},
									},
									{
										Step: &pb.Pipeline_Step{
											Name:      "other",
											DependsOn: []string{"prep"},
											Kind: &pb.Pipeline_Step_Exec_{
												Exec: &pb.Pipeline_Step_Exec{
													Image:   "busybox",
													Command: "echo",
													Args:    []string{"other"},
												},
											},
										},
										Application: appRef,
										Job:         &pb.Ref_Job{Id: "job-for-other-step"},
										State:       pb.UI_PipelineRunTreeNode_RUNNING,
										StartTime:   quickTimestamp("2023-01-01T13:00:50Z"),
										Children: &pb.UI_PipelineRunTreeNode_Children{
											Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
											Nodes: []*pb.UI_PipelineRunTreeNode{},
										},
									},
								},
							},
						},
						{
							Step: &pb.Pipeline_Step{
								Name:      "done",
								DependsOn: []string{"invoke-referenced-pipeline"},
								Kind: &pb.Pipeline_Step_Exec_{
									Exec: &pb.Pipeline_Step_Exec{
										Image:   "busybox",
										Command: "echo",
										Args:    []string{"done"},
									},
								},
							},
							Application: appRef,
							Job:         &pb.Ref_Job{Id: "job-for-done-step"},
							State:       pb.UI_PipelineRunTreeNode_QUEUED,
							Children: &pb.UI_PipelineRunTreeNode_Children{
								Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
								Nodes: []*pb.UI_PipelineRunTreeNode{},
							},
						},
					},
				},
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			result, err := UI_PipelineRunTreeFromJobs(tt.Jobs, tt.StatusReports)

			require.NoError(err)

			if diff := cmp.Diff(tt.Tree, result, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected difference:\n%v", diff)
			}

		})
	}
}

// quickTimestamp parses an RFC3339-formatted string and returns the time it
// represents as a timestamppb.Timestamp.
//
// This is intended purely to make tests more readable and robust to daylight
// savings time etc.
func quickTimestamp(s string) *timestamppb.Timestamp {
	t, _ := time.Parse(time.RFC3339, s)
	return timestamppb.New(t)
}
