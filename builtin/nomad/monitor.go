package nomad

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
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

// NewMonitor returns a new monitor. The returned monitor will
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
func (m *monitor) Monitor(ctx context.Context, evalID string) error {
	// Add the initial pending state
	m.update(newEvalState())
	eval, _, err := m.client.Evaluations().Info(evalID, nil)

	stream := m.client.EventStream()
	topics := make(map[api.Topic][]string)
	topicName := api.Topic("Evaluation:" + eval.JobID)
	topics["Evaluation"] = append(topics[topicName], evalID)
	eventStream, err := stream.Stream(ctx, topics, 0, nil)
	if err != nil {
		return err
	}

	// We wait for 5 minutes for the evaluation to finish
	d := time.Now().Add(time.Minute * time.Duration(5))
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
		case <-ctx.Done(): // cancelled
			return status.Errorf(codes.Aborted, "Context cancelled: %s", ctx.Err())
		}
		events := <-eventStream
		for _, event := range events.Events {
			if event.Topic != "Evaluation" {
				return status.Errorf(codes.FailedPrecondition, "Evaluations API did not return evaluation")
			}
			if event.Payload["Evaluation"] != nil {
				evalPayload := event.Payload["Evaluation"].(map[string]interface{})
				if evalPayload["JobID"].(string) == eval.JobID {
					switch evalPayload["Status"] {
					case "pending":
						continue
					case "failed":
						return status.Errorf(codes.FailedPrecondition, "Evaluation failed")
					case "complete":
						return nil
					case "cancelled":
						return status.Errorf(codes.FailedPrecondition, "Evaluation cancelled")
					case "blocked":
						// We error here because even though Nomad can start a new eval if a job can only be partially deployed
						// (where the new eval is "blocked" until allocs can be scheduled, we error here because we only
						// check the initial eval ID (for now)
						return status.Errorf(codes.FailedPrecondition, "Evaluation blocked")
					default:
						return status.Errorf(codes.FailedPrecondition, "Unknown evaluation status: %s", evalPayload["Status"].(string))
					}
				} else {
					// TODO: Log warning - we shouldn't get evaluations for jobs we didn't want
					break
				}
			}
		}
	}
}
