package telemetry

import (
	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"github.com/hashicorp/go-hclog"
	ocview "go.opencensus.io/stats/view"
	octrace "go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// datadogExporter implements the telemetry exporter interface.
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
