// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package runner

// cleanup stacks cleanup functions to call when Close is called.
func (r *Runner) cleanup(f func()) {
	oldF := r.cleanupFunc
	r.cleanupFunc = func() {
		defer f()
		if oldF != nil {
			oldF()
		}
	}
}
