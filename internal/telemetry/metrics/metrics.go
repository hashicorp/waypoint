// Package metrics is a wrapper around the OpenTelemetry libraries providing
// support for emitting metrics. It implements a very simple usage pattern,
// which can typically express the logging of a metric value in a single
// line of code, reducing the boilerplate required.
package metrics

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	// "go.opentelemetry.io/otel/metric"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/metric"
	// "go.opentelemetry.io/otel/unit"
)

// Metrics is a wrapper around OpenTelemetry's metrics library.
type Metrics struct {
	parent *Metrics
	// meter     metric.Meter
	prefix    string
	durations map[string]*stats.Float64Measure
	// samples   map[string]metric.Float64ValueRecorder
	// gauges    map[string]*gauge
	options []Option

	l sync.RWMutex

	// Configured exporters that need to be registered and closed
	// TODO: (clint) telemetry exporter needs to be here, or datadog options?
}

// New returns a new metrics instance.
func New(prefix string) *Metrics {
	return &Metrics{
		// meter:     m,
		prefix:    prefix,
		durations: make(map[string]*stats.Float64Measure),
		// samples:   make(map[string]metric.Float64ValueRecorder),
		// gauges:    make(map[string]*gauge),
	}
}

// NewChild creates a new metrics sub-instance, applying the prefix and
// copying all current attributes.
func (m *Metrics) NewChild(prefix string) *Metrics {
	if m == nil {
		return New(prefix)
	}

	return &Metrics{
		parent:    m,
		prefix:    prefix,
		durations: make(map[string]*stats.Float64Measure),
	}
}

// AddDuration is used to add durations for metrics
func (m *Metrics) AddDuration(name string, stat *stats.Float64Measure) {
	if m == nil {
		log.Println("&&&&&&&&&&&&&&&")
		log.Println("&&&&&&&& M IS NILL not adding =======")
		log.Println("&&&&&&&&&&&&&&&")
		return
	} else {
		log.Println("&&&&&&&&&&&&&&&")
		log.Println("&&&&&&&& M set, adding.. =======")
		log.Println("&&&&&&&&&&&&&&&")
	}

	// TODO: load first/error handle
	m.l.Lock()
	m.durations[name] = stat
	l := len(m.durations)
	log.Println("===============")
	log.Println(fmt.Sprintf("======== duration len in Add (%d) =======", l))
	log.Println("===============")
	m.l.Unlock()
}

// SetAttribute is used to set an attribute key/value pair on the metrics
// instance, which will be logged with each metric thereafter.
func (m *Metrics) SetAttribute(key, value string) {
	if m == nil {
		return
	}

	// Skip attributes with missing key or value.
	if key == "" || value == "" {
		return
	}

	m.l.Lock()
	m.options = append(m.options, WithAttribute(key, value))
	m.l.Unlock()
}

// StartTimer returns a new *Timer, which can be used to measure the elapsed
// time since StartTimer was called.
func (m *Metrics) StartTimer(name string, options ...Option) *Timer {
	if m == nil {
		return noop().StartTimer(name, options...)
	}

	return &Timer{
		metrics: m,
		name:    name,
		start:   time.Now(),
		options: options,
	}
}

// MeasureSince is used to measure the time elapsed since t.
func (m *Metrics) MeasureSince(ctx context.Context, name string, t time.Time, options ...Option) {
	if m == nil {
		log.Println("&&&&&&&&&&&&&&&")
		log.Println("&&&&&&&& M IS NILL in measure since =======")
		log.Println("&&&&&&&&&&&&&&&")
		return
	} else {
		log.Println("&&&&&&&&&&&&&&&")
		log.Println("&&&&&&&& M is set, measuring.... =======")
		log.Println("&&&&&&&&&&&&&&&")
	}
	if m == nil {
		return
	}

	// if m.parent != nil {
	// 	// m.parent.MeasureSince(m.metricName(name), t, m.mergeOptions(options)...)
	// 	m.parent.MeasureSince(m.metricName(name), t, m.mergeOptions(options)...)
	// 	return
	// }

	m.l.RLock()
	r, ok := m.durations[Jobs]
	l := len(m.durations)
	m.l.RUnlock()

	if !ok {
		// var err error

		log.Println("===============")
		log.Println(fmt.Sprintf("======== err jobs not found (%d) =======", l))
		log.Println("===============")
		return
		// r, err = m.meter.NewInt64ValueRecorder(m.metricName(name)+".milliseconds",
		// 	metric.WithUnit(unit.Milliseconds))
		// if err != nil {
		// 	return
		// }
		// r = stats.Int64("repl2/latency", "Other distribution", "ms")

		// m.l.Lock()
		// m.durations[name] = r
		// m.l.Unlock()
	}

	// opts := m.evalOptions(options...)

	// r.Record(opts.Context, time.Since(t).Milliseconds(), opts.attrs()...)
	// _ = stats.RecordWithTags(context.Background(), []tag.Mutator{
	// tag.Upsert(clusterIDTag, clusterID),
	// }, measurement.M(1))
	ts := sinceInMilliseconds(t)
	log.Printf("=== Debug sending record with tags")
	err := stats.RecordWithTags(context.Background(), []tag.Mutator{
		tag.Upsert(KeyJobType, m.prefix),
		tag.Upsert(KeyServerVersion, "0.8.1"),
		// }, r.M(time.Since(t).Milliseconds()))
	}, r.M(ts))

	log.Printf("=== time since: %v", ts)
	log.Printf("=== err value record with tags: %v", err)
	// stats.Record(context.Background(), r.M(time.Since(t).Milliseconds()))
}

// // SetGauge is used to set the value of a named gauge. The value will be logged
// // continuously until later modified. Note that the given options, if any, are
// // only evaluated on the first call for the named gauge. Subsequent calls do
// // not respect any updated options. value must be a builtin numerical type. All
// // other types will be discarded.
// func (m *Metrics) SetGauge(name string, value interface{}, options ...Option) {
// 	if m == nil {
// 		return
// 	}

// 	floatValue, ok := float64Value(value)
// 	if !ok {
// 		return
// 	}

// 	if m.parent != nil {
// 		m.parent.SetGauge(m.metricName(name), value,
// 			m.mergeOptions(options)...)
// 		return
// 	}

// 	m.l.RLock()
// 	g, ok := m.gauges[name]
// 	m.l.RUnlock()

// 	if !ok {
// 		g = &gauge{val: floatValue}

// 		opts := m.evalOptions(options...)

// 		var err error
// 		g.observer, err = m.meter.NewFloat64ValueObserver(
// 			m.metricName(name),
// 			func(_ context.Context, r metric.Float64ObserverResult) {
// 				r.Observe(g.value(), opts.attrs()...)
// 			},
// 		)
// 		if err != nil {
// 			return
// 		}

// 		m.l.Lock()
// 		m.gauges[name] = g
// 		m.l.Unlock()
// 	}

// 	g.setValue(floatValue)
// }

// RecordValue records the given value as a single, point-in-time sample.
// value must be builtin numerical type. Other types will be discarded.
func (m *Metrics) RecordValue(name string, value interface{}, options ...Option) {
	if m == nil {
		return
	}
}

// 	floatValue, ok := float64Value(value)
// 	if !ok {
// 		return
// 	}

// 	if m.parent != nil {
// 		m.parent.RecordValue(m.metricName(name), value,
// 			m.mergeOptions(options)...)
// 		return
// 	}

// 	m.l.RLock()
// 	r, ok := m.samples[name]
// 	m.l.RUnlock()

// 	if !ok {
// 		var err error

// 		r, err = m.meter.NewFloat64ValueRecorder(m.metricName(name))
// 		if err != nil {
// 			return
// 		}

// 		m.l.Lock()
// 		m.samples[name] = r
// 		m.l.Unlock()
// 	}

// 	opts := m.evalOptions(options...)

// 	r.Record(opts.Context, floatValue, opts.attrs()...)
// }

// // evalOptions is used to apply the given options in the order given.
// func (m *Metrics) evalOptions(opts ...Option) *MetricOptions {
// 	options := &MetricOptions{
// 		Context:    context.Background(),
// 		Attributes: make(map[string]string),
// 	}

// 	for _, opt := range m.mergeOptions(opts) {
// 		opt(options)
// 	}

// 	return options
// }

// metricName applies any configured prefix and returns the full metric name.
func (m *Metrics) metricName(name string) string {
	if m.prefix != "" {
		return m.prefix + "." + name
	}
	return name
}

// mergeOptions merges the instance-level options with the metric-level
// options given, producing a fully contextualized set of Options.
func (m *Metrics) mergeOptions(opts []Option) []Option {
	m.l.RLock()
	defer m.l.RUnlock()

	return append(m.options, opts...)
}

// Timer holds information about a duration metric to be logged in the future,
// and provides a convenient interface for recording it.
type Timer struct {
	metrics *Metrics
	name    string
	start   time.Time
	options []Option
}

// Record records the current elapsed time since the start of the Timer.
func (t *Timer) Record(ctx context.Context) {
	t.metrics.MeasureSince(ctx, t.name, t.start, t.options...)
}

// MetricOptions carries information about options given at the call site
// when logging metrics.
type MetricOptions struct {
	Context    context.Context
	Attributes map[string]string
}

// // attrs converts the map-style key/value attributes for use with the
// // OpenTelemetry APIs.
// func (m *MetricOptions) attrs() []attribute.KeyValue {
// 	attrs := make([]attribute.KeyValue, 0, len(m.Attributes))
// 	for k, v := range m.Attributes {
// 		attrs = append(attrs, attribute.String(k, v))
// 	}
// 	return attrs
// }

// Option represents an option passed in at the call site when logging metrics.
type Option func(*MetricOptions)

// // WithContext produces an Option which carries the given context through to
// // the underlying OpenTelemetry API calls when logging metrics.
// func WithContext(ctx context.Context) Option {
// 	return func(m *MetricOptions) {
// 		m.Context = ctx
// 	}
// }

// WithAttributes produces an Option which carries the given attribute through
// to the underlying OpenTelemetry API calls when logging metrics.
func WithAttribute(key, value string) Option {
	return func(m *MetricOptions) {
		m.Attributes[key] = value
	}
}

// // gauge is the internal representation of a gauge and its value.
// type gauge struct {
// 	observer metric.Float64ValueObserver
// 	val      float64

// 	sync.RWMutex
// }

// // setValue safely sets the value of the gauge. The value will be recorded
// // during the next observer run.
// func (g *gauge) setValue(value float64) {
// 	g.Lock()
// 	g.val = value
// 	g.Unlock()
// }

// // value safely returns the current value of the gauge.
// func (g *gauge) value() float64 {
// 	g.RLock()
// 	defer g.RUnlock()

// 	return g.val
// }

// noop returns a metrics instance that doesn't do anything.
func noop() *Metrics {
	return New("")
}

// // float64Value converts the given value to a float64.
// func float64Value(value interface{}) (float64, bool) {
// 	switch v := value.(type) {
// 	case float64:
// 		return v, true
// 	case float32:
// 		return float64(v), true
// 	case uint64:
// 		return float64(v), true
// 	case uint32:
// 		return float64(v), true
// 	case uint16:
// 		return float64(v), true
// 	case uint8:
// 		return float64(v), true
// 	case int64:
// 		return float64(v), true
// 	case int32:
// 		return float64(v), true
// 	case int16:
// 		return float64(v), true
// 	case int8:
// 		return float64(v), true
// 	case int:
// 		return float64(v), true
// 	default:
// 		return 0, false
// 	}
// }

func sinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}
