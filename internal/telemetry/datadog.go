// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package telemetry

import (
	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"github.com/hashicorp/go-hclog"
	ocview "go.opencensus.io/stats/view"
	octrace "go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// datadogExporter is a wrapper around a datadog.Exporter that implements exporter to give some well-defined registration
// and shut down behavior.
type datadogExporter struct {
	exporter *datadog.Exporter
	log      hclog.Logger
}

func (d *datadogExporter) register() {
	d.log.Debug("Registering the Datadog exporter")
	octrace.RegisterExporter(d.exporter)
	ocview.RegisterExporter(d.exporter)
}

func (d *datadogExporter) close() error {
	d.log.Debug("Shutting down Datadog exporter")
	d.exporter.Stop()
	d.log.Debug("Datadog exporter stopped")
	return nil
}

func newDatadogExporter(opts datadog.Options, log hclog.Logger) (*datadogExporter, error) {
	ddExporter, err := datadog.NewExporter(opts)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to initialize Datadog exporter: %s", err)
	}
	return &datadogExporter{
		exporter: ddExporter,
		log:      log,
	}, nil
}
