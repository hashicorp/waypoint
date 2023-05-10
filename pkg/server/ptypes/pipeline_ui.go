// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ptypes

import (
	"fmt"
	"sort"

	"github.com/hashicorp/waypoint/internal/pkg/graph"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

// UI_PipelineRunTreeFromJobs takes a graph of Jobs and StatusReports and
// returns a specialized tree representation that drives the pipeline-run
// timeline UI.
func UI_PipelineRunTreeFromJobs(jobs []*pb.Job, statusReports []*pb.StatusReport) (*pb.UI_PipelineRunTreeNode, error) {
	processor := newUIPipelineProcessor(jobs, statusReports)

	return processor.run()
}

// newUIPipelineProcessor creates a new uiPipelineProcessor for the given jobs and status reports
func newUIPipelineProcessor(jobs []*pb.Job, statusReports []*pb.StatusReport) uiPipelineProcessor {
	p := uiPipelineProcessor{
		jobIdx:  make(map[string]*pb.Job),
		nodeIdx: make(map[string]*pb.UI_PipelineRunTreeNode),
	}

	// Populate job index
	for _, job := range jobs {
		p.jobIdx[job.Id] = job
	}

	// Populate status report index
	statusReportIdx := make(map[string]*pb.StatusReport)
	for _, statusReport := range statusReports {
		if d := statusReport.GetDeploymentId(); d != "" {
			statusReportIdx[d] = statusReport
		} else if r := statusReport.GetReleaseId(); r != "" {
			statusReportIdx[r] = statusReport
		}
	}

	// Populate node index
	for _, job := range jobs {
		var step *pb.Pipeline_Step
		var statusReport *pb.StatusReport

		if stepOp := job.GetPipelineStep(); stepOp != nil {
			step = stepOp.Step
		} else {
			// This job represents a nested or referenced pipeline
			// invocation. It doesn’t encode all the necessary
			// information to infer the step details, so we’ll
			// create a placeholder step and fill in the details as
			// we learn more from other jobs in the set.
			step = &pb.Pipeline_Step{
				Name:      job.Pipeline.Step,
				DependsOn: []string{},
				Kind:      &pb.Pipeline_Step_Pipeline_{},
			}
		}

		// Look up the latest status report for this node
		if d := job.Result.GetDeploy().GetDeployment(); d != nil {
			statusReport = statusReportIdx[d.Id]
		}
		if r := job.Result.GetRelease().GetRelease(); r != nil {
			statusReport = statusReportIdx[r.Id]
		}

		// Build the node (children will be populated later)
		node := &pb.UI_PipelineRunTreeNode{
			Step:               step,
			Job:                &pb.Ref_Job{Id: job.Id},
			State:              p.nodeStateFromJob(job),
			StartTime:          job.AckTime,
			CompleteTime:       job.CompleteTime,
			Application:        job.Application,
			Workspace:          job.Workspace,
			Result:             job.Result,
			LatestStatusReport: statusReport,
			Children: &pb.UI_PipelineRunTreeNode_Children{
				Mode:  pb.UI_PipelineRunTreeNode_Children_SERIAL,
				Nodes: []*pb.UI_PipelineRunTreeNode{},
			},
		}

		// Register the node the index
		p.nodeIdx[job.Id] = node
	}

	// Build graph from job `DependsOn` fields
	for i, job := range p.jobIdx {
		p.graph.Add(i)

		for _, j := range job.DependsOn {
			if _, exists := p.jobIdx[j]; !exists {
				// If the referenced job isn’t in the set then
				// we can safely ignore it. It’s likely a
				// task-releated ancillary job.
				continue
			}

			p.graph.Add(j)
			p.graph.AddEdge(j, i)
		}
	}

	return p
}

// uiPipelineProcessor encapsulates all the state required to transform a set of
// jobs and status reports into our desired output
type uiPipelineProcessor struct {
	jobIdx  map[string]*pb.Job
	nodeIdx map[string]*pb.UI_PipelineRunTreeNode
	graph   graph.Graph
}

// run performs our battery of transformations and returns the root of the
// output tree, or an error if something went wrong.
func (p uiPipelineProcessor) run() (*pb.UI_PipelineRunTreeNode, error) {
	var rootId string
	var rootNode *pb.UI_PipelineRunTreeNode

	// Remove shortcuts from the graph. We do this because sometimes jobs
	// contain redundant entries in their DependsOn list.
	p.graph.TransitiveReduction()

	// Find the root id/node
	rootId, rootNode, err := p.root()
	if err != nil {
		return nil, err
	}

	// Process the tree starting at the root ID in the context of the root node
	if err := p.processSubGraph(rootId, rootNode); err != nil {
		return nil, err
	}

	return rootNode, nil
}

// root returns the root ID and node.
func (p uiPipelineProcessor) root() (string, *pb.UI_PipelineRunTreeNode, error) {
	sorted := p.graph.KahnSort()

	id, ok := sorted[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("could not find root node")
	}

	node, ok := p.nodeIdx[id]
	if !ok {
		return "", nil, fmt.Errorf("could not find root node")
	}

	return id, node, nil
}

// processSubGraph performs transformations starting at the given inputId, and
// collecting the results into the given outputNode.
func (p uiPipelineProcessor) processSubGraph(inputId string, outputNode *pb.UI_PipelineRunTreeNode) error {
	var nextIds []string
	for _, v := range p.graph.OutEdges(inputId) {
		nextIds = append(nextIds, v.(string))
	}
	sort.Strings(nextIds)
	degree := len(nextIds)

	if degree == 1 {
		nextId := nextIds[0]
		nextNode := p.nodeIdx[nextId]
		inputPipelineId := p.jobIdx[inputId].Pipeline.GetPipelineId()
		nextPipelineId := p.jobIdx[nextId].Pipeline.GetPipelineId()

		if inputPipelineId != nextPipelineId {
			invokeId, err := p.foldSubPipeline(nextId, inputPipelineId)
			if err != nil {
				return err
			}
			nextId = invokeId
			nextNode = p.nodeIdx[nextId]
		}

		outputNode.Children.Nodes = append(outputNode.Children.Nodes, nextNode)

		return p.processSubGraph(nextId, outputNode)
	}

	if degree > 1 {
		virtualId, virtualNode, err := p.foldBranches(inputId, outputNode)
		if err != nil {
			return err
		}

		// Continue processing the parent graph from the new
		// virtual node, but accumulating into the output node.
		if err := p.processSubGraph(virtualId, outputNode); err != nil {
			return err
		}

		// Process all the input children as sub-graphs,
		// accumulating into the virtual node.
		for _, nextId := range nextIds {
			nextNode := p.nodeIdx[nextId]
			inputJob := p.jobIdx[inputId]
			nextJob := p.jobIdx[nextId]
			inputPipelineId := inputJob.Pipeline.GetPipelineId()
			nextPipelineId := nextJob.Pipeline.GetPipelineId()

			if inputPipelineId != nextPipelineId {
				invokeId, err := p.foldSubPipeline(nextId, inputPipelineId)
				if err != nil {
					return err
				}
				nextId = invokeId
				nextNode = p.nodeIdx[nextId]
			}

			virtualNode.Children.Nodes = append(virtualNode.Children.Nodes, nextNode)

			if err := p.processSubGraph(nextId, nextNode); err != nil {
				return err
			}
		}

		p.inferNodeAttrsFromChildren(virtualNode)
	}

	return nil
}

// foldBranches replaces parallel branches of a graph with a single node.
//
// Before:
//
//	  A
//	┌─┴─┐
//	B   D
//	C   E
//	└─┬─┘
//	  F
//
// After:
//
//	A
//	V
//	F
//
// It returns the new node and assumes the parent function will perform further
// processing to embed the original children within the new node.
func (p uiPipelineProcessor) foldBranches(
	inputId string,
	outputNode *pb.UI_PipelineRunTreeNode,
) (string, *pb.UI_PipelineRunTreeNode, error) {
	nextIds := p.graph.OutEdges(inputId)
	degree := len(nextIds)

	// Create a “virtual node” to encapsulate concurrent work.
	virtualId := outputNode.Job.Id + "-virtual"
	virtualNode := &pb.UI_PipelineRunTreeNode{
		Job: &pb.Ref_Job{
			Id: virtualId,
		},
		Children: &pb.UI_PipelineRunTreeNode_Children{
			Mode:  pb.UI_PipelineRunTreeNode_Children_PARALLEL,
			Nodes: []*pb.UI_PipelineRunTreeNode{},
		},
	}
	virtualJob := &pb.Job{
		Pipeline: p.jobIdx[inputId].Pipeline,
	}
	p.nodeIdx[virtualId] = virtualNode
	p.jobIdx[virtualId] = virtualJob

	// Add the virtual node into the graph as a child of the
	// input node.
	p.graph.Add(virtualId)
	p.graph.AddEdge(inputId, virtualId)
	// And add the virtual node to the output node’s children.
	outputNode.Children.Nodes = append(outputNode.Children.Nodes, virtualNode)

	// Disconnect all children from the input node.
	for _, nextId := range nextIds {
		p.graph.RemoveEdge(inputId, nextId)
	}

	// Find common descendent (if any)
	var commonDescendent string
	seen := make(map[graph.Vertex]int)
	for _, i := range nextIds {
		err := p.graph.DFS(i, func(j graph.Vertex, c func() error) error {
			seen[j] += 1
			if seen[j] == degree {
				commonDescendent = j.(string)
				return nil
			} else {
				return c()
			}
		})
		if err != nil {
			return "", nil, err
		}
	}

	if commonDescendent != "" {
		// If there is a common ancestor, disconnect
		// anything that’s pointing to it.
		for _, a := range p.graph.InEdges(commonDescendent) {
			p.graph.RemoveEdge(a, commonDescendent)
		}
		// And connect the virtual node instead.
		p.graph.AddEdge(virtualId, commonDescendent)
	}

	return virtualId, virtualNode, nil
}

// foldSubPipeline replaces nested pipeline invocations with a single node.
func (p uiPipelineProcessor) foldSubPipeline(inputId string, parentPipelineId string) (string, error) {
	var invokeId string
	inputNode := p.nodeIdx[inputId]

	// Find the node in which we return to the parent pipeline. This is the
	// `use "pipeline"` step from the pipeline.
	err := p.graph.DFS(inputId, func(v graph.Vertex, c func() error) error {
		id := v.(string)
		job := p.jobIdx[id]

		if job.Pipeline.GetPipelineId() == parentPipelineId {
			invokeId = id
			return nil
		}

		return c()
	})
	if err != nil {
		return "", err
	}
	if invokeId == "" {
		subPipelineId := p.jobIdx[inputId].Pipeline.GetPipelineId()
		return "", fmt.Errorf("Invocation node not found for sub-pipeline %q", subPipelineId)
	}

	// Disconnect the invoke node from it’s current inbound edges
	for _, i := range p.graph.InEdges(invokeId) {
		p.graph.RemoveEdge(i, invokeId)
	}

	// Moves inbound edges from the input node to the invoke node
	for _, i := range p.graph.InEdges(inputId) {
		p.graph.RemoveEdge(i, inputId)
		p.graph.AddEdge(i, invokeId)
	}

	invokeNode := p.nodeIdx[invokeId]

	// Embed the original input node in the invoke node
	invokeNode.Children.Nodes = append(invokeNode.Children.Nodes, inputNode)

	// Process the sub-graph starting at the input node
	if err := p.processSubGraph(inputId, invokeNode); err != nil {
		return "", err
	}

	// Extract sub-pipeline information from children
	invokeNode.Step.Kind = p.stepRefFromJob(p.jobIdx[inputId])
	p.inferNodeAttrsFromChildren(invokeNode)

	return invokeId, nil
}

// nodeStateFromJob return the “tree node state” for a given job. The tree node
// state is somewhat abstracted/inferred from job state, thus the need for this
// mapping.
func (p uiPipelineProcessor) nodeStateFromJob(j *pb.Job) pb.UI_PipelineRunTreeNode_State {
	switch j.State {
	case pb.Job_QUEUED, pb.Job_WAITING:
		return pb.UI_PipelineRunTreeNode_QUEUED
	case pb.Job_RUNNING:
		return pb.UI_PipelineRunTreeNode_RUNNING
	case pb.Job_ERROR:
		if j.CancelTime != nil {
			return pb.UI_PipelineRunTreeNode_CANCELLED
		} else {
			return pb.UI_PipelineRunTreeNode_ERROR
		}
	case pb.Job_SUCCESS:
		return pb.UI_PipelineRunTreeNode_SUCCESS
	}

	return pb.UI_PipelineRunTreeNode_UNKNOWN
}

// stepRefFromJob returns a pipeline step definition from a given job. This
// function doesn’t do anything terribly complex, really it just encapsulates a
// lot of protobuf boilerplate.
func (p uiPipelineProcessor) stepRefFromJob(job *pb.Job) *pb.Pipeline_Step_Pipeline_ {
	return &pb.Pipeline_Step_Pipeline_{
		Pipeline: &pb.Pipeline_Step_Pipeline{
			Ref: &pb.Ref_Pipeline{
				Ref: &pb.Ref_Pipeline_Owner{
					Owner: &pb.Ref_PipelineOwner{
						Project:      &pb.Ref_Project{Project: job.Application.Project},
						PipelineName: job.Pipeline.PipelineName,
					},
				},
			},
		},
	}
}

// inferNodeAttrsFromChildren takes a node and sets the following attributes
// based on the attributes of its children:
//
// * Application
// * StartTime
// * CompleteTime
// * State
func (p uiPipelineProcessor) inferNodeAttrsFromChildren(node *pb.UI_PipelineRunTreeNode) {
	children := node.Children.Nodes

	if len(children) == 0 {
		return
	}

	// Extract application
	node.Application = children[0].Application

	// Extract start time
	for _, n := range children {
		t1 := node.StartTime
		t2 := n.StartTime

		if t1 == nil {
			node.StartTime = t2
			continue
		}

		if t2 == nil {
			continue
		}

		if t2.AsTime().Before(t1.AsTime()) {
			node.StartTime = t2
		}
	}

	// Extract complete time
	for _, n := range node.Children.Nodes {
		t1 := node.CompleteTime
		t2 := n.CompleteTime

		if t1 == nil {
			node.CompleteTime = t2
			continue
		}

		if t2 == nil {
			// We’ve found an incomplete step, so the parent is also incomplete
			node.CompleteTime = nil
			break
		}

		if t2.AsTime().After(t1.AsTime()) {
			node.CompleteTime = t2
		}
	}

	// Extract state
	for _, n := range node.Children.Nodes {
		s := n.State
		node.State = s

		if s != pb.UI_PipelineRunTreeNode_SUCCESS {
			break
		}
	}
}
