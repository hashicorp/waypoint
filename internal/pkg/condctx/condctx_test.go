// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package condctx

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	// Create cond var and wait on it
	cond := sync.NewCond(&sync.Mutex{})
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		cond.L.Lock()
		defer cond.L.Unlock()
		cond.Wait()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	closer := Notify(ctx, cond)
	defer closer()

	// We should still be waiting
	select {
	case <-doneCh:
		t.Fatal("should wait")

	case <-time.After(100 * time.Millisecond):
		// Good
	}

	// Cancel
	cancel()

	// We should be woken up
	select {
	case <-doneCh:

	case <-time.After(500 * time.Millisecond):
		t.Fatal("should wake up")
	}

	// Cancel again does nothing
	cancel()
}
