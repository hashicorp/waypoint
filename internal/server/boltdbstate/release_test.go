// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"testing"
)

func TestRelease(t *testing.T) {
	releaseOp.Test(t)
}
