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

type Telemetry struct {
	config struct {
		enableOpenCensusExporter  bool
		openCensusExporterOptions []ocagent.ExporterOption

		enableDatadogExporter  bool
		datadogExporterOptions datadog.Options

		enableZpages bool
		// Address to listen on for zpages. Defaults to 127.0.0.1:9999
		zpagesAddr string
	}

	// Logger is the logger to use. This will default to hclog.L() if not set.
	log hclog.Logger

	// Configured and running exporters that need to be closed
	zpagesServer *http.Server

	// Configured exporters that need to be registered and closed
	exporters []Exporter
}

// Exporter is an OpenCensus exporter
type Exporter interface {
	Register()
	Close() error
}

// exporter is a prototype struct that can be embedded to anonymously fulfill the Exporter interface.
type exporter struct {
	register func()
	close    func() error
}

func (e *exporter) Register()    { e.register() }
func (e *exporter) Close() error { return e.close() }

// NewTelemetry initializes the telemetry components.
func NewTelemetry(opts ...Option) (Telemetry, error) {
	var t Telemetry
	for _, opt := range opts {
		opt(&t)
	}

	if t.log == nil {
		t.log = hclog.L().Named("telemetry")
	}
	log := t.log

	config := t.config

	if config.enableOpenCensusExporter {
		log.Debug("Creating the OpenCensus agent exporter")

		e, err := openCensusAgentExporter(config.openCensusExporterOptions, log)
		if err != nil {
			return t, err
		}

		t.exporters = append(t.exporters, e)
	}

	if config.enableDatadogExporter {
		log.Debug("Starting the Datadog exporter")

		e, err := datadogExporter(config.datadogExporterOptions, log)
		if err != nil {
			return t, err
		}

		t.exporters = append(t.exporters, e)
	}

	// Run zPages
	if config.enableZpages {
		zPagesMux := http.NewServeMux()
		zpages.Handle(zPagesMux, "/debug")

		var zpagesAddr string
		if config.zpagesAddr == "" {
			zpagesAddr = "127.0.0.1:9999"
		} else {
			zpagesAddr = config.zpagesAddr
		}

		srv := http.Server{
			Addr:    zpagesAddr,
			Handler: zPagesMux,
		}

		t.zpagesServer = &srv
	}

	// Less frequent sampling can be achieved by exporting to an OpenCensus collector with sampling configured.
	octrace.ApplyConfig(octrace.Config{DefaultSampler: octrace.AlwaysSample()})

	return t, nil
}

func openCensusAgentExporter(opts []ocagent.ExporterOption, log hclog.Logger) (Exporter, error) {
	ocExporter, err := ocagent.NewExporter(opts...)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to initialize OpenCensus agent exporter: %s", err)
	}

	// Create our exporter that we can register and close later
	e := &struct{ exporter }{}

	e.register = func() {
		log.Debug("Registering the OpenCensus agent exporter")
		octrace.RegisterExporter(ocExporter)
		ocview.RegisterExporter(ocExporter)
	}

	e.close = func() error {
		log.Debug("Shutting down OpenCensus agent exporter")
		ocExporter.Flush()
		if err := ocExporter.Stop(); err != nil {
			return fmt.Errorf("failed to stop the OpenCensus agent exporter: %s", err)
		} else {
			log.Debug("OpenCensus agent exporter flushed and stopped")
		}
		return nil
	}
	return e, nil
}

func datadogExporter(opts datadog.Options, log hclog.Logger) (Exporter, error) {
	ddExporter, err := datadog.NewExporter(opts)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to initialize Datadog exporter: %s", err)
	}

	// Create our exporter that we can register and close later
	e := &struct{ exporter }{}

	e.register = func() {
		log.Debug("Registering the Datadog exporter")
		octrace.RegisterExporter(ddExporter)
		ocview.RegisterExporter(ddExporter)
	}

	e.close = func() error {
		log.Debug("Shutting down Datadog exporter")
		ddExporter.Stop()
		log.Debug("Datadog exporter flushed and stopped")
		return nil
	}

	return e, nil
}

// Run registers and starts the telemetry providers. It blocks until the provided context closes.
func (t *Telemetry) Run(ctx context.Context) error {
	log := t.log

	// Register all of our configured exporters
	for _, e := range t.exporters {
		e.Register()
	}

	// Run zPages
	if t.zpagesServer != nil {
		go func() {
			log.Debug("Starting zPages server", "addr", t.zpagesServer.Addr)
			if err := t.zpagesServer.ListenAndServe(); err != nil {
				log.Debug("zPages server exited", "err", err)
			}
		}()
	}

	// Wait on context to close
	<-ctx.Done()

	// Close zPages
	if t.zpagesServer != nil {
		log.Debug("Shutting down zPages server")
		err := t.zpagesServer.Close()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to shut down zPages server: %s", err))
		}
	}

	log.Debug("Shutting down telemetry exporters")
	for _, e := range t.exporters {
		if err := e.Close(); err != nil {
			log.Error("Failed to close exporter", "err", err)
		}
	}

	log.Debug("Finished shutting down telemetry components")
	return nil
}

type Option func(*Telemetry)

// WithLogger sets the logger for use with the server.
func WithLogger(log hclog.Logger) Option {
	return func(t *Telemetry) {
		t.log = log
	}
}

func WithOpenCensusExporter(exporterOptions []ocagent.ExporterOption) Option {
	return func(t *Telemetry) {
		t.config.enableOpenCensusExporter = true
		t.config.openCensusExporterOptions = exporterOptions
	}
}

func WithDatadogExporter(exporterOptions datadog.Options) Option {
	return func(t *Telemetry) {
		t.config.enableDatadogExporter = true
		t.config.datadogExporterOptions = exporterOptions
	}
}

// WithZpages enables a zpages server. Addr will default to 127.0.0.1:9999 unless otherwise specified.
func WithZpages(addr string) Option {
	return func(t *Telemetry) {
		t.config.enableZpages = true
		t.config.zpagesAddr = addr
	}
}
