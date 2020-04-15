package component

import (
	"context"
	"encoding/base32"
	"time"

	"golang.org/x/crypto/blake2b"
)

// LogPlatform is responsible for reading the logs for a deployment.
// This doesn't need to be the same as the Platform but a Platform can also
// implement this interface to natively provide logs.
type LogPlatform interface {
	// LogsFunc should return an implementation of LogViewer.
	LogsFunc() interface{}
}

// LogViewer returns batches of log lines. This is expected to be returned
// by a LogPlatform implementation.
type LogViewer interface {
	// NextBatch is called to return the next batch of logs. This is expected
	// to block if there are no logs available. The context passed in will be
	// cancelled if the logs viewer is interrupted.
	NextLogBatch(ctx context.Context) ([]LogEvent, error)
}

// LogEvent represents a single log entry.
type LogEvent struct {
	Partition string
	Timestamp time.Time
	Message   string
}

type PartitionViewer struct {
	shortened map[string]string
}

var encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

func (pv *PartitionViewer) Short(part string) string {
	if len(part) < 10 {
		return part
	}

	if pv.shortened == nil {
		pv.shortened = make(map[string]string)
	}

	if short, ok := pv.shortened[part]; ok {
		return short
	}

	h, _ := blake2b.New(blake2b.Size, nil)

	h.Write([]byte(part))

	str := encoding.EncodeToString(h.Sum(nil))

	length := 7

	for {
		short := str[:length]
		if _, found := pv.shortened[short]; found {
			length++
		} else {
			pv.shortened[part] = short
			return short
		}
	}
}
