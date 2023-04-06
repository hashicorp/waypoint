// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"testing"
)

func TestRelease(t *testing.T) {
	releaseOp.Test(t)
}
