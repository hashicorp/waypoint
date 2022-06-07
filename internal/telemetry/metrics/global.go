package metrics

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"go.opencensus.io/stats"
	// "go.opentelemetry.io/otel/metric"
)

// The global metrics instance, which leverages the global OpenTelemetry
// metrics exporter.
var global atomic.Value

func init() {
	// Store a typed nil. The implementation supports nil receivers.
	log.Println("===============")
	log.Println("==== init call ===")
	log.Println("===============")
	global.Store((*Metrics)(nil))
}

// NewGlobal creates a new global metrics instance.
func NewGlobal(prefix string) {
	global.Store(New(prefix))
}

// NewChild creates a metrics sub-instance based on the global instance.
func NewChild(prefix string) *Metrics {
	return global.Load().(*Metrics).NewChild(prefix)
}

// SetAttribute sets the given key/value pair on the global metrics instance.
func SetAttribute(key, value string) {
	global.Load().(*Metrics).SetAttribute(key, value)
}

// StartTimer uses the global instance to return a new *Timer, which can be
// used to measure the elapsed time since StartTimer was called.
func StartTimer(name string, options ...Option) *Timer {
	return global.Load().(*Metrics).StartTimer(name, options...)
}

// AddDuration uses the global instance to add durations
func AddDuration(name string, stat *stats.Float64Measure) {
	global.Load().(*Metrics).AddDuration(name, stat)
}

// MeasureSince measures the time elapsed since t using the global instance.
func MeasureSince(ctx context.Context, name string, t time.Time, options ...Option) {
	global.Load().(*Metrics).MeasureSince(ctx, name, t, options...)
}

// // SetGauge sets the value of the named gauge using the global instance.
// func SetGauge(name string, value interface{}, options ...Option) {
// 	global.Load().(*Metrics).SetGauge(name, value, options...)
// }

// // RecordValue records the given float64 sample using the global instance.
// func RecordValue(name string, value interface{}, options ...Option) {
// 	global.Load().(*Metrics).RecordValue(name, value, options...)
// }

// // EmitRuntimeMetrics is a blocking routine which periodically emits runtime
// // metrics using the global instance.
// func EmitRuntimeMetrics(ctx context.Context, interval time.Duration) {
// 	global.Load().(*Metrics).EmitRuntimeMetrics(ctx, interval)
// }
