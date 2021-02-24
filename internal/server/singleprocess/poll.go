package singleprocess

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// runPollQueuer starts the poll queuer. The poll queuer sleeps on and
// schedules polling operations for projects that have polling enabled.
// This blocks and is expected to be run in a goroutine.
//
// This function should only ever be invoked one at a time. Running multiple
// copies can result in duplicate polls for projects.
func (s *service) runPollQueuer(
	ctx context.Context,
	wg *sync.WaitGroup,
	funclog hclog.Logger,
) {
	defer wg.Done()

	funclog.Info("starting")
	defer funclog.Info("exiting")

	for {
		log := funclog

		if ctx.Err() != nil {
			// If our context was cancelled, exit.
			return
		}

		ws := memdb.NewWatchSet()
		p, pollTime, err := s.state.ProjectPollPeek(ws)
		if err != nil {
			// This error really should never happen. Instead of just exiting,
			// we log it and just sleep a minute. Hopefully someone will notice
			// the logs. We sleep for a minute because any error that happened
			// here is probably real bad and is gonna keep happening.
			log.Error("BUG (please report): error during poll queuer, sleeping 1 minute", "err", err)
			time.Sleep(1 * time.Minute)
			continue
		}
		if p != nil {
			log = log.With("project", p.Name)
		}

		var loopCtxCancel context.CancelFunc
		loopCtx := ctx
		if !pollTime.IsZero() {
			loopCtx, loopCtxCancel = context.WithDeadline(ctx, pollTime)
		}

		// Confusing bit below. Here is the explanation of the problem we're
		// solving for: there are THREE possible outcomes that we are waiting on:
		//
		//   (1) WatchSet (ws) triggers - this means that the data changed,
		//       i.e. a project changed polling settings, so we need to reloop.
		//
		//   (2) ctx is cancelled - this means the whole queuer is cancelled
		//       and we just want to exit.
		//
		//   (3) loopCtx is cancelled - this means we hit our deadline for
		//       polling and we want to queue a polling operation for this
		//       project.
		//

		log.Trace("waiting on watchset and contexts")
		err = ws.WatchCtx(loopCtx)
		loopCtxErr := loopCtx.Err()
		if loopCtxCancel != nil {
			loopCtxCancel()
		}

		if err == nil {
			// Outcome (1) above
			log.Debug("dataset change, restarting poll queuer")
			continue
		}

		if ctx.Err() != nil {
			// Outcome (2) above
			return
		}

		if loopCtxErr == nil {
			// Should never happen since by here we should definitely
			// be in outcome (3) but if this happened then... just restart
			// cause its weird.
			log.Warn("poll deadline wasn't hit but watchset triggered")
			continue
		}

		// p is allowed to be nil in this loop, but it should never reach
		// this point. Given we use it below, we put this check here to warn
		// loudly that it happened. p shouldn't be nil here because if p is
		// nil then we have no pollTime and therefore no loopCtx either. This
		// means outcome (1) or (2) MUST happen.
		if p == nil {
			log.Error("reached outcome (3) in poller with nil p. This should not happen.")
			continue
		}

		// Outcome (3)
		log.Trace("queueing poll job")
		resp, err := s.QueueJob(ctx, &pb.QueueJobRequest{
			Job: &pb.Job{
				// SingletonId so that we only have one poll operation at
				// any time queued per project.
				SingletonId: fmt.Sprintf("poll/%s", p.Name),

				Application: &pb.Ref_Application{
					Project: p.Name,
					// No Application set since PollOp is project-oriented
				},

				// Polling always happens on the default workspace even
				// though the PollOp is across every workspace.
				Workspace: &pb.Ref_Workspace{Workspace: "default"},

				// Poll!
				Operation: &pb.Job_Poll{
					Poll: &pb.Job_PollOp{},
				},

				// Any runner is fine for polling.
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Any{
						Any: &pb.Ref_RunnerAny{},
					},
				},
			},
		})
		if err != nil {
			log.Warn("error queueing a poll job", "err", err)
			continue
		}
		log.Debug("queued polling job", "job_id", resp.JobId)

		// Mark this as complete so the next poll gets rescheduled.
		log.Trace("scheduling next project poll time")
		if err := s.state.ProjectPollComplete(p, time.Now()); err != nil {
			// This should never happen so like above, if this happens we
			// sleep for a minute so we don't completely overload the
			// server since this is likely to happen again. We want people
			// to see this in the logs.
			log.Warn("BUG (please report): error marking project polling complete", "err", err)
			time.Sleep(1 * time.Minute)
			continue
		}
	}
}
