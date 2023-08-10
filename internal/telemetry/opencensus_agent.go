// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package telemetry

import (
	"fmt"

	"contrib.go.opencensus.io/exporter/ocagent"
	"github.com/hashicorp/go-hclog"
	ocview "go.opencensus.io/stats/view"
	octrace "go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// opencensusAgentExporter is a wrapper around an ocagent.Exporter that implements exporter to give some well-defined registration
// and shut down behavior.
type opencensusAgentExporter struct {
	exporter *ocagent.Exporter
	log      hclog.Logger
}

func newOpenCensusAgentExporter(opts []ocagent.ExporterOption, log hclog.Logger) (*opencensusAgentExporter, error) {
	ocExporter, err := ocagent.NewExporter(opts...)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to initialize OpenCensus agent exporter: %s", err)
	}
	return &opencensusAgentExporter{
		exporter: ocExporter,
		log:      log,
	}, nil
}

func (o *opencensusAgentExporter) register() {
	o.log.Debug("Registering the OpenCensus agent exporter")
	octrace.RegisterExporter(o.exporter)
	ocview.RegisterExporter(o.exporter)
}

func (o *opencensusAgentExporter) close() error {
	o.log.Debug("Shutting down OpenCensus agent exporter")
	o.exporter.Flush()
	if err := o.exporter.Stop(); err != nil {
		return fmt.Errorf("failed to stop the OpenCensus agent exporter: %s", err)
	} else {
		o.log.Debug("OpenCensus agent exporter flushed and stopped")
	}
	return nil
}
