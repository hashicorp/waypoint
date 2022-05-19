package nomad

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func (m *monitor) Monitor(ctx context.Context, evalID string) error {
	// Add the initial pending state
	m.update(newEvalState())
	eval, _, err := m.client.Evaluations().Info(evalID, nil)

	events := m.client.EventStream()
	topics := make(map[api.Topic][]string)
	topics["Evaluation"] = append(topics["Evaluation"], evalID)
	eventStream, err := events.Stream(ctx, topics, 0, nil)
	if err != nil {
		return err
	}

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
		event := <-eventStream
		for _, event := range event.Events {
			if event.Topic != "Evaluation" {
				return errors.New("Evaluations API did not return evaluation")
			}
			if event.Payload["Evaluation"] != nil {
				evalPayload := event.Payload["Evaluation"].(map[string]string)
				if evalPayload["JobID"] == eval.JobID {
					switch evalPayload["Status"] {
					case "pending":
						continue
					case "failed":
						errors.New("Evaluation failed")
					case "completed":
						return nil
					case "cancelled":
						return errors.New("Evaluation cancelled")
					default:
						return errors.New(fmt.Sprintf("Unknown evaluation status: %s", evalPayload["Status"]))
					}
				} else {
					// TODO: Log warning - we shouldn't get evaluations for jobs we didn't want
					break
				}

			}
		}
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
