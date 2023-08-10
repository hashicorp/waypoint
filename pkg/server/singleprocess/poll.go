// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package singleprocess

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

//go:generate mockery -case underscore -structname PollHandler -name pollHandler

// pollHandler is a private interface that the server implements for polling on
// different items such as projects or status reports.
type pollHandler interface {
	// Peek returns the next item that should be polled.
	// This will return (nil,nil,nil) if there are no items to poll currently.
	//
	// This calls the items state implementation of its "peek" operation so it
	// does not update the poll item's next poll time. Therefore, calling this
	// multiple times should return the same result unless a function like
	// Complete is called.
	//
	// Note that the WatchSet must be populated with a watch channel that is triggered when
	// there might be a new or changed record
	Peek(hclog.Logger, memdb.WatchSet) (interface{}, time.Time, error)

	// PollJob generates a QueueJobRequest that is used to poll on.
	// It is expected to be given a proto message obtained from Peek which
	// is used to define the job returned.
	PollJob(hclog.Logger, interface{}) ([]*pb.QueueJobRequest, error)

	// Complete will mark the job that was queued as complete using the specific
	// state implementation.
	Complete(hclog.Logger, interface{}) error
}

// runPollQueuer starts the poll queuer. The poll queuer sleeps on and
// schedules polling operations for pollable items that have polling enabled
// and implemented.
// This blocks and is expected to be run in a goroutine.
//
// This function should only ever be invoked one at a time. Running multiple
// copies can result in duplicate polls for items.
func (s *Service) runPollQueuer(
	ctx context.Context,
	wg *sync.WaitGroup,
	handler pollHandler,
	funclog hclog.Logger,
) {
	defer wg.Done()

	// We allow nil loggers cause most funcs that take an hclog do.
	// Use default logger in this case.
	if funclog == nil {
		funclog = hclog.L()
	}

	funclog.Info("starting")
	defer funclog.Info("exiting")

	for {
		log := funclog

		if ctx.Err() != nil {
			// If our context was cancelled, exit.
			return
		}

		ws := memdb.NewWatchSet()
		pollItem, pollTime, err := handler.Peek(log, ws)
		if err != nil {
			// This error really should never happen. Instead of just exiting,
			// we log it and just sleep a minute. Hopefully someone will notice
			// the logs. We sleep for a minute because any error that happened
			// here is probably real bad and is gonna keep happening.
			log.Error("BUG (please report): error during poll queuer, sleeping 1 minute", "err", err)

			// We also exit on context done so we can just exit the goroutine.
			select {
			case <-time.After(1 * time.Minute):
			case <-ctx.Done():
			}

			continue
		}

		var loopCtxCancel context.CancelFunc
		loopCtx := ctx
		if !pollTime.IsZero() {
			log.Trace("next poll time", "time", pollTime)
			loopCtx, loopCtxCancel = context.WithDeadline(ctx, pollTime)
		}

		// Confusing bit below. Here is the explanation of the problem we're
		// solving for: there are THREE possible outcomes that we are waiting on:
		//
		//   (1) WatchSet (ws) triggers - this means that the data changed,
		//       i.e. a poll item changed polling settings, so we need to reloop.
		//
		//   (2) ctx is cancelled - this means the whole queuer is cancelled
		//       and we just want to exit.
		//
		//   (3) loopCtx is cancelled - this means we hit our deadline for
		//       polling and we want to queue a polling operation for this
		//       poll item.
		//

		log.Trace("waiting on watchset and contexts")
		err = ws.WatchCtx(loopCtx)
		loopCtxErr := loopCtx.Err()
		if loopCtxCancel != nil {
			loopCtxCancel()
		}

		if err == nil {
			// Outcome (1) above
			log.Trace("dataset change, restarting poll queuer")
			continue
		}

		if ctx.Err() != nil {
			// Outcome (2) above
			log.Trace("context cancelled for poll queuer, returning from poll loop ctx")
			return
		}

		if loopCtxErr == nil {
			// Should never happen since by here we should definitely
			// be in outcome (3) but if this happened then... just restart
			// cause its weird.
			log.Warn("poll deadline wasn't hit but watchset triggered")
			continue
		}

		// pollItem is allowed to be nil in this loop, but it should never reach
		// this point. Given we use it below, we put this check here to warn
		// loudly that it happened. pollItem shouldn't be nil here because if
		// pollItem is nil then we have no pollTime and therefore no loopCtx either.
		// This means outcome (1) or (2) MUST happen.
		if pollItem == nil {
			log.Error("reached outcome (3) in poller with nil pollItem. " +
				"This should not happen. This usually means there is a bug " +
				"in the pollHandler implementation")
			time.Sleep(1 * time.Second)
			continue
		}

		// Outcome (3)
		log.Trace("queueing poll jobs")
		queueJobRequests, err := handler.PollJob(log, pollItem)
		if err != nil {
			log.Warn("error building a poll job request. This should not happen "+
				"repeatedly. If you see this in your log repeatedly, report a bug.",
				"err", err)
			time.Sleep(1 * time.Second)
			continue
		}

		totalRequests := len(queueJobRequests)
		log.Trace("queueing jobs for poller", "job_total", totalRequests)

		// Note: We queue all poll jobs transactionally and return any
		// errors that occurred
		_, err = s.queueJobMulti(ctx, queueJobRequests)
		if err != nil {
			log.Warn("error queueing a poll job", "err", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Mark this as complete so the next poll gets rescheduled.
		log.Trace("completing poll and scheduling next poll time")
		if err := handler.Complete(log, pollItem); err != nil {
			// This should never happen so like above, if this happens we
			// sleep for a minute so we don't completely overload the
			// server since this is likely to happen again. We want people
			// to see this in the logs.
			log.Warn("BUG (please report): error marking polling item complete", "err", err)
			time.Sleep(1 * time.Minute)
			continue
		}
	}
}
