// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package boltdbstate

import (
	"testing"
)

func TestStatusReport(t *testing.T) {
	statusReportOp.Test(t)
}
