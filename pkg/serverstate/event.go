// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serverstate

import (
	"time"

	"github.com/hashicorp/waypoint/pkg/server/gen"
)

type Event struct {
	Application    *gen.Ref_Application
	Project        *gen.Project
	EventTimestamp time.Time
	EventType      string
	EventData      []byte
}
