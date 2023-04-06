// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

const (
	// updateWait is the amount of time to wait between status
	// updates. Because the monitor is poll-based, we use this
	// delay to avoid overwhelming the API server.
	updateWait = time.Second
)

// evalState is used to store the current "state of the world"
// in the context of monitoring an evaluation.
type evalState struct {
	status     string
	desc       string
	node       string
	deployment string
	job        string
	allocs     map[string]*allocState
	wait       time.Duration
	index      uint64
}

// newEvalState creates and initializes a new monitorState
func newEvalState() *evalState {
	return &evalState{
		status: "pending",
		allocs: make(map[string]*allocState),
	}
}

// allocState is used to track the state of an allocation
type allocState struct {
	id          string
	group       string
	node        string
	desired     string
	desiredDesc string
	client      string
	clientDesc  string
	index       uint64
}

// monitor wraps an evaluation monitor and holds metadata and
// state information.
type monitor struct {
	ui     terminal.Status
	client *api.Client
	state  *evalState

	sync.Mutex
}

// newMonitor returns a new monitor. The returned monitor will
// write output information to the provided ui.
func NewMonitor(ui terminal.Status, client *api.Client) *monitor {
	mon := &monitor{
		ui:     ui,
		client: client,
		state:  newEvalState(),
	}
	return mon
}

// update is used to update our monitor with new state. It can be
// called whether the passed information is new or not, and will
// only dump update messages when state changes.
func (m *monitor) update(update *evalState) {
	m.Lock()
	defer m.Unlock()

	existing := m.state

	// Swap in the new state at the end
	defer func() {
		m.state = update
	}()

	// Check the allocations
	for allocID, alloc := range update.allocs {
		if existing, ok := existing.allocs[allocID]; !ok {
			switch {
			case alloc.index < update.index:
				// New alloc with create index lower than the eval
				// create index indicates modification
				m.ui.Step(terminal.StatusOK, fmt.Sprintf(
					"Allocation %q modified: node %q, group %q",
					alloc.id, alloc.node, alloc.group))

			case alloc.desired == "run":
				// New allocation with desired status running
				m.ui.Step(terminal.StatusOK, fmt.Sprintf(
					"Allocation %q created: node %q, group %q",
					alloc.id, alloc.node, alloc.group))
			}
		} else {
			switch {
			case existing.client != alloc.client:
				description := ""
				if alloc.clientDesc != "" {
					description = fmt.Sprintf(" (%s)", alloc.clientDesc)
				}
				// Allocation status has changed
				m.ui.Step(terminal.StatusOK, fmt.Sprintf(
					"Allocation %q status changed: %q -> %q%s",
					alloc.id, existing.client, alloc.client, description))
			}
		}
	}

	// Check if the status changed. We skip any transitions to pending status.
	if existing.status != "" &&
		update.status != "pending" &&
		existing.status != update.status {
		m.ui.Step(terminal.StatusOK, fmt.Sprintf("Evaluation status changed: %q -> %q",
			existing.status, update.status))
	}
}

// Monitor is used to start monitoring the given evaluation ID. It
// writes output directly to the terminal ui, and returns an error.
func (m *monitor) Monitor(evalID string) error {
	// Track if we encounter a scheduling failure. This can only be
	// detected while querying allocations, so we use this bool to
	// carry that status into the return code.
	var schedFailure bool

	// Add the initial pending state
	m.update(newEvalState())

	for {
		// Query the evaluation
		eval, _, err := m.client.Evaluations().Info(evalID, nil)
		if err != nil {
			return fmt.Errorf("No evaluation with id %q found", evalID)
		}
		m.ui.Update(fmt.Sprintf("Monitoring evaluation %q", eval.ID))

		// Create the new eval state.
		state := newEvalState()
		state.status = eval.Status
		state.desc = eval.StatusDescription
		state.node = eval.NodeID
		state.job = eval.JobID
		state.deployment = eval.DeploymentID
		state.wait = eval.Wait
		state.index = eval.CreateIndex

		// Query the allocations associated with the evaluation
		allocs, _, err := m.client.Evaluations().Allocations(eval.ID, nil)
		if err != nil {
			return fmt.Errorf("Error reading allocations: %s", err)
		}

		// Add the allocs to the state
		for _, alloc := range allocs {
			state.allocs[alloc.ID] = &allocState{
				id:          alloc.ID,
				group:       alloc.TaskGroup,
				node:        alloc.NodeID,
				desired:     alloc.DesiredStatus,
				desiredDesc: alloc.DesiredDescription,
				client:      alloc.ClientStatus,
				clientDesc:  alloc.ClientDescription,
				index:       alloc.CreateIndex,
			}
		}

		// Update the state
		m.update(state)

		switch eval.Status {
		case "complete", "failed", "cancelled":
			if len(eval.FailedTGAllocs) == 0 {
				m.ui.Step(terminal.StatusOK, fmt.Sprintf("Evaluation %q finished with status %q",
					eval.ID, eval.Status))
			} else {
				// There were failures making the allocations
				schedFailure = true
				m.ui.Step(terminal.StatusWarn, fmt.Sprintf("Evaluation %q finished with status %q but failed to place all allocations",
					eval.ID, eval.Status))

				// Print the failures per task group
				for tg, metrics := range eval.FailedTGAllocs {
					noun := "allocation"
					if metrics.CoalescedFailures > 0 {
						noun += "s"
					}
					output := fmt.Sprintf("Task Group %q (failed to place %d %s):\n", tg, metrics.CoalescedFailures+1, noun)
					metrics := formatAllocMetrics(metrics, false, "  ")
					for _, line := range strings.Split(metrics, "\n") {
						output = output + line
					}

					m.ui.Step(terminal.StatusWarn, output)
				}

				if eval.BlockedEval != "" {
					m.ui.Step(terminal.StatusWarn, fmt.Sprintf("Evaluation %q waiting for additional capacity to place remainder",
						eval.BlockedEval))
				}
			}
		default:
			// Wait for the next update
			time.Sleep(updateWait)
			continue
		}

		break
	}
	if schedFailure {
		return fmt.Errorf("Failed to schedule all allocations")
	}

	return nil
}

func formatAllocMetrics(metrics *api.AllocationMetric, scores bool, prefix string) string {
	// Print a helpful message if we have an eligibility problem
	var out string
	if metrics.NodesEvaluated == 0 {
		out += fmt.Sprintf("%s* No nodes were eligible for evaluation\n", prefix)
	}

	// Print a helpful message if the user has asked for a DC that has no
	// available nodes.
	for dc, available := range metrics.NodesAvailable {
		if available == 0 {
			out += fmt.Sprintf("%s* No nodes are available in datacenter %q\n", prefix, dc)
		}
	}

	// Print filter info
	for class, num := range metrics.ClassFiltered {
		out += fmt.Sprintf("%s* Class %q: %d nodes excluded by filter\n", prefix, class, num)
	}
	for cs, num := range metrics.ConstraintFiltered {
		out += fmt.Sprintf("%s* Constraint %q: %d nodes excluded by filter\n", prefix, cs, num)
	}

	// Print exhaustion info
	if ne := metrics.NodesExhausted; ne > 0 {
		out += fmt.Sprintf("%s* Resources exhausted on %d nodes\n", prefix, ne)
	}
	for class, num := range metrics.ClassExhausted {
		out += fmt.Sprintf("%s* Class %q exhausted on %d nodes\n", prefix, class, num)
	}
	for dim, num := range metrics.DimensionExhausted {
		out += fmt.Sprintf("%s* Dimension %q exhausted on %d nodes\n", prefix, dim, num)
	}

	// Print quota info
	for _, dim := range metrics.QuotaExhausted {
		out += fmt.Sprintf("%s* Quota limit hit %q\n", prefix, dim)
	}

	out = strings.TrimSuffix(out, "\n")
	return out
}
