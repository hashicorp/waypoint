package serverstate

import (
	"github.com/hashicorp/waypoint/pkg/server/gen"
	"time"
)

type Event struct {
	Application    gen.Ref_Application
	EventTimestamp time.Time
	EventType      string
	EventData      []byte
}