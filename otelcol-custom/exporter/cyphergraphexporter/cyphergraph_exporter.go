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

func newTracesExporter(cfg *Config, set exporter.CreateSettings) (*cyphergraphTraceExporter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	var neo4jAuth neo4j.AuthToken
	if cfg.Username != "" {
		neo4jAuth = neo4j.BasicAuth(cfg.Username, string(cfg.Password), "")
	} else if t := string(cfg.BearerToken); t != "" {
		neo4jAuth = neo4j.BearerAuth(t)
	} else if t := string(cfg.KerberosTicket); t != "" {
		neo4jAuth = neo4j.KerberosAuth(t)
	} else {
		neo4jAuth = neo4j.NoAuth()
	}
	driver, err := neo4j.NewDriverWithContext(
		cfg.DatabaseURI,
		neo4jAuth,
		withLogger(set.Logger),
		withUserAgent(cfg.UserAgent),
		withBackOffConfig(&cfg.BackOffConfig),
	)
	if err != nil {
		return nil, err
	}
	return &cyphergraphTraceExporter{
		driver: driver,
		logger: set.Logger,
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
