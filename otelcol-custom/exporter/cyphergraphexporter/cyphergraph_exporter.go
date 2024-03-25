package cyphergraphexporter // import "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type cyphergraphTraceExporter struct {
	driver neo4j.DriverWithContext
	logger *zap.Logger
}

func newTracesExporter(config *Config, settings exporter.CreateSettings) (*cyphergraphTraceExporter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	driver, err := neo4j.NewDriverWithContext(
		config.DatabaseURI,
		neo4j.BasicAuth(config.Username, string(config.Password), ""),
		withLogger(settings.Logger),
		withUserAgent(config.UserAgent),
		withBackOffConfig(&config.BackOffConfig),
	)
	if err != nil {
		return nil, err
	}
	return &cyphergraphTraceExporter{
		driver: driver,
		logger: settings.Logger,
	}, nil
}

func (e *cyphergraphTraceExporter) Start(ctx context.Context, _ component.Host) error {
	err := e.driver.VerifyConnectivity(ctx)
	if err != nil {
		return err
	}
	// TODO think about initializing database indices here.
	return nil
}

func (e *cyphergraphTraceExporter) Shutdown(ctx context.Context) error {
	if e.driver == nil {
		return nil
	}
	return e.driver.Close(ctx)
}

func (e *cyphergraphTraceExporter) tracesPusher(ctx context.Context, td ptrace.Traces) error {
	// TODO implement
	_, _ = ctx, td
	return fmt.Errorf("not implemented")
}
