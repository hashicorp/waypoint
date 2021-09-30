package telemetry

import (
	"context"
	"net/http"

	"go.opencensus.io/zpages"

	ocview "go.opencensus.io/stats/view"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"contrib.go.opencensus.io/exporter/ocagent"

	octrace "go.opencensus.io/trace"

	"github.com/hashicorp/go-hclog"
)

type Option func(*telemetry)

func Run(opts ...Option) error {
	var t telemetry
	for _, opt := range opts {
		opt(&t)
	}

	// TODO(izaak): would be nice to warn if no exporters have been configured

	var closeFuncs []func()
	if t.EnableOpenCensusExporter {
		exporter, err := ocagent.NewExporter(t.OpenCensusExporterOptions...)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "failed to initalize opencensus agent exporter: %s", err)
		}
		octrace.RegisterExporter(exporter)
		ocview.RegisterExporter(exporter)
		closeFuncs = append(closeFuncs, func() {
			t.Log.Debug("Shutting down opencensus agent exporter")
			exporter.Flush()
			if err := exporter.Stop(); err != nil {
				t.Log.Error("Failed to stop the opencensus agent exporter: %s", err)
			} else {
				t.Log.Debug("OpenCensus agent exporter flushed and stopped")
			}
		})
	}

	// Run Zpages
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
			t.Log.Debug("Starting zPages server at %s", zpagesAddr)
			if err := srv.ListenAndServe(); err != nil {
				panic("Failed to serve zPages")
			}
		}()
		closeFuncs = append(closeFuncs, func() {
			err := srv.Close()
			if err != nil {
				t.Log.Error("Failed to shut down zPages server: %s", err)
			}
		})
	}

	// TODO: allow applying sampling different sampling config
	octrace.ApplyConfig(octrace.Config{DefaultSampler: octrace.AlwaysSample()})

	// Wait on context to close
	select {
	case <-t.Context.Done():
		for _, f := range closeFuncs {
			f()
		}
		t.Log.Debug("Finished shutting down telemetry components")
	}
	return nil
}

type telemetry struct {
	// Context is the context to use for the telemetry agents. When this is cancelled,
	// the agents will gracefully shut down.
	Context context.Context

	// Logger is the logger to use. This will default to hclog.L() if not set.
	Log hclog.Logger

	// TODO(izaak): looks like ocagent doesn't do global tags? This would be nice...
	// Tags to apply globally if possible
	//GlobalTags []string

	EnableOpenCensusExporter  bool
	OpenCensusExporterOptions []ocagent.ExporterOption

	EnableZpages bool
	// Address to listen on for zpages. Defaults to 127.0.0.1:9999
	ZpagesAddr string

	// TODO(izaak): datadog exporter options
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
		t.Log = log
	}
}

func WithOpenCensusExporter(exporterOptions []ocagent.ExporterOption) Option {
	return func(t *telemetry) {
		t.EnableOpenCensusExporter = true
		t.OpenCensusExporterOptions = exporterOptions
	}
}

// WithZpages enables a zpages server. Addr will default to 127.0.0.1:9999 unless otherwise specified.
func WithZpages(addr string) Option {
	return func(t *telemetry) {
		t.EnableZpages = true
		t.ZpagesAddr = addr
	}
}
