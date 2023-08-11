// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package serverstate

import (
	"time"

	"github.com/hashicorp/waypoint/pkg/server/gen"
)

type Event struct {
	Application    *gen.Ref_Application
	Project        *gen.Ref_Project
	EventTimestamp time.Time
	EventType      string
	EventData      []byte
}
