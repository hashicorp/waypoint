// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"testing"
)

func TestStatusReport(t *testing.T) {
	statusReportOp.Test(t)
}
