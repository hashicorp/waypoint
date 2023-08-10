// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package funcs

import (
	"testing"
	"time"

	"github.com/zclconf/go-cty/cty"
)

func TestTimestamp(t *testing.T) {
	currentTime := time.Now().UTC()

	result, err := TimestampFunc.Call([]cty.Value{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	resultTime, err := time.Parse(time.RFC3339, result.AsString())
	if err != nil {
		t.Fatalf("Error parsing timestamp: %s", err)
	}

	if resultTime.Sub(currentTime).Seconds() > 10.0 {
		t.Fatalf("Timestamp Diff too large. Expected: %s\nReceived: %s", currentTime.Format(time.RFC3339), result.AsString())
	}
}
