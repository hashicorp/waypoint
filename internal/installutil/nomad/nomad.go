package nomad

import (
	"context"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"time"
)

func RunJob(
	ctx context.Context,
	s terminal.Step,
	client *api.Client,
	job *api.Job,
	policyOverride bool,
) (string, error) {
	jobOpts := &api.RegisterOptions{
		PolicyOverride: policyOverride,
	}

	resp, _, err := client.Jobs().RegisterOpts(job, jobOpts, nil)
	if err != nil {
		return "", err
	}

	s.Update("Waiting for allocation to be scheduled")
	qopts := &api.QueryOptions{
		WaitIndex: resp.EvalCreateIndex,
		WaitTime:  time.Duration(500 * time.Millisecond),
	}

	eval, meta, err := waitForEvaluation(ctx, s, client, resp, qopts)
	if err != nil {
		return "", err
	}
	if eval == nil {
		return "", fmt.Errorf("evaluation status could not be determined")
	}
	qopts.WaitIndex = meta.LastIndex

	var allocID string
	retries := 0
	maxRetries := 3
	for {
		allocs, qmeta, err := client.Evaluations().Allocations(eval.ID, qopts)
		if err != nil {
			return "", err
		}
		qopts.WaitIndex = qmeta.LastIndex
		if len(allocs) == 0 {
			return "", fmt.Errorf("no allocations found after evaluation completed")
		}

		switch allocs[0].ClientStatus {
		case "running":
			allocID = allocs[0].ID
			s.Update("Nomad allocation running")
			retries++
		case "pending":
			s.Update(fmt.Sprintf("Waiting for allocation %q to start", allocs[0].ID))
			// retry
		default:
			return "", fmt.Errorf("allocation failed")
		}

		if allocID != "" {
			if retries == maxRetries {
				return allocID, nil
			} else {
				s.Update("Ensuring allocation %q has properly started up...", allocs[0].ID)
				time.Sleep(1 * time.Second)
			}
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	return "", nil
}

func waitForEvaluation(
	ctx context.Context,
	s terminal.Step,
	client *api.Client,
	resp *api.JobRegisterResponse,
	qopts *api.QueryOptions,
) (*api.Evaluation, *api.QueryMeta, error) {

	for {
		eval, meta, err := client.Evaluations().Info(resp.EvalID, qopts)
		if err != nil {
			return nil, nil, err
		}

		qopts.WaitIndex = meta.LastIndex

		switch eval.Status {
		case "pending":
			s.Update("Nomad allocation pending...")
		case "complete":
			s.Update("Nomad allocation created")

			return eval, meta, nil
		case "failed", "canceled", "blocked":
			s.Update("Nomad failed to schedule the job")
			s.Status(terminal.StatusError)
			return nil, nil, fmt.Errorf("Nomad evaluation did not transition to 'complete'")
		default:
			return nil, nil, fmt.Errorf("receieved unknown eval status from Nomad: %q", eval.Status)
		}
	}
}
