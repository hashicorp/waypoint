// Package condctx provides helpers for working with condition variables
// along with the standard "context" package and interface.
package condctx

import (
	"context"
	"sync"
)

// Notify will wake up all waiters of cond when the context is cancelled.
// To use this, callers should call Notify and then check ctx.Err in their
// condition loop for the condition variable.
//
// The return value should be deferred to properly clean up resources associated
// with this function.
func Notify(ctx context.Context, cond *sync.Cond) func() {
	doneCh := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			// Wake up all condition vars so we wake ourself up.
			cond.Broadcast()

		case <-doneCh:
			// Return since we were cancelled.
		}
	}()

	return func() {
		// We do this if so that the function can be called multiple times.
		if doneCh != nil {
			close(doneCh)
			doneCh = nil
		}
	}
}
