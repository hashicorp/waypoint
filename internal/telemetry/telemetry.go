package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"contrib.go.opencensus.io/exporter/ocagent"
	datadog "github.com/DataDog/opencensus-go-exporter-datadog"
	"github.com/hashicorp/go-hclog"
	ocview "go.opencensus.io/stats/view"
	octrace "go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Option func(*telemetry)

func Run(opts ...Option) error {
	var t telemetry
	for _, opt := range opts {
		opt(&t)
	}

	if t.Logger == nil {
		t.Logger = hclog.L().Named("telemetry")
	}
	log := t.Logger

	// If any of the below agents or servers need to close on exit, they can
	// register their close function here, and they will be called when the parent
	// context closes this runner.
	var closeFuncs []func()

	if t.EnableOpenCensusExporter {
		log.Debug("Starting the opencensus agent exporter")
		exporter, err := ocagent.NewExporter(t.OpenCensusExporterOptions...)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to initalize opencensus agent exporter: %s", err)
		}
		octrace.RegisterExporter(exporter)
		ocview.RegisterExporter(exporter)
		closeFuncs = append(closeFuncs, func() {
			log.Debug("Shutting down OpenCensus agent exporter")
			exporter.Flush()
			if err := exporter.Stop(); err != nil {
				log.Error("Failed to stop the opencensus agent exporter", "err", err)
			} else {
				log.Debug("OpenCensus agent exporter flushed and stopped")
			}
		})
	}

	if t.EnableDatadogExporter {
		log.Debug("Starting the opencensus datadog exporter")
		exporter, err := datadog.NewExporter(t.DatadogExporterOptions)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to initalize datadog exporter: %s", err)
		}
		octrace.RegisterExporter(exporter)
		ocview.RegisterExporter(exporter)
		closeFuncs = append(closeFuncs, func() {
			log.Debug("Shutting down datadog exporter")
			exporter.Stop()
			log.Debug("Datadog exporter flushed and stopped")
		})
	}

	// Run zPages
	if t.EnableZpages {
		zPagesMux := http.NewServeMux()
		zpages.Handle(zPagesMux, "/debug")

		var zpagesAddr string
		if t.ZpagesAddr == "" {
			zpagesAddr = "127.0.0.1:9999"
		} else {
			zpagesAddr = t.ZpagesAddr
		}

		srv := http.Server{
			Addr:    zpagesAddr,
			Handler: zPagesMux,
		}

		go func() {
			log.Debug("Starting zPages server", "addr", zpagesAddr)
			if err := srv.ListenAndServe(); err != nil {
				log.Debug("zPages server exited", "err", err)
			}
		}()
		closeFuncs = append(closeFuncs, func() {
			err := srv.Close()
			if err != nil {
				log.Error(fmt.Sprintf("Failed to shut down zPages server: %s", err))
			}
		})
	}

	// Less frequent sampling can be achieved by exporting to an opencensus collector with sampling configured.
	octrace.ApplyConfig(octrace.Config{DefaultSampler: octrace.AlwaysSample()})

	// Wait on context to close
	<-t.Context.Done()
	log.Debug("Shutting down telemetry components")
	for _, f := range closeFuncs {
		f()
	}
	log.Debug("Finished shutting down telemetry components")
	return nil
}

type telemetry struct {
	// Context is the context to use for the telemetry agents. When this is cancelled,
	// the agents will gracefully shut down.
	Context context.Context

	// Logger is the logger to use. This will default to hclog.L() if not set.
	Logger hclog.Logger

	EnableOpenCensusExporter  bool
	OpenCensusExporterOptions []ocagent.ExporterOption

	EnableDatadogExporter  bool
	DatadogExporterOptions datadog.Options

	EnableZpages bool
	// Address to listen on for zpages. Defaults to 127.0.0.1:9999
	ZpagesAddr string
}

// WithContext sets the logger for use with the server.
func WithContext(ctx context.Context) Option {
	return func(t *telemetry) {
		t.Context = ctx
	}
}

// WithLogger sets the logger for use with the server.
func WithLogger(log hclog.Logger) Option {
	return func(t *telemetry) {
		t.Logger = log
	}
}

func WithOpenCensusExporter(exporterOptions []ocagent.ExporterOption) Option {
	return func(t *telemetry) {
		t.EnableOpenCensusExporter = true
		t.OpenCensusExporterOptions = exporterOptions
	}
}

func WithDatadogExporter(exporterOptions datadog.Options) Option {
	return func(t *telemetry) {
		t.EnableDatadogExporter = true
		t.DatadogExporterOptions = exporterOptions
	}
}

// WithZpages enables a zpages server. Addr will default to 127.0.0.1:9999 unless otherwise specified.
func WithZpages(addr string) Option {
	return func(t *telemetry) {
		t.EnableZpages = true
		t.ZpagesAddr = addr
	}
}
