package cyphergraphexporter // import "github.com/clement-casse/Playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type cyphergraphTraceExporter struct{}

func (e *cyphergraphTraceExporter) tracesPusher(ctx context.Context, td ptrace.Traces) error {
	// TODO implement
	return fmt.Errorf("not implemented")
}

func newTracesExporter(config Config, settings exporter.CreateSettings) (*cyphergraphTraceExporter, error) {
	return nil, fmt.Errorf("not implemented")
}
