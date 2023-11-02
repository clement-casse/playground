package cyphergraphexporter // import "github.com/clement-casse/Playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/clement-casse/Playground/otelcol-custom/exporter/cyphergraphexporter/internal/metadata"
)

const (
	defaultUsername    = ""
	defaultPassword    = ""
	defaultDatabaseUri = "bolt://localhost:7687"
)

// NewFactory creates a factory for the CypherGraph exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, metadata.TracesStability))
}

func createDefaultConfig() component.Config {
	return &Config{
		Username:    defaultUsername,
		Password:    defaultPassword,
		DatabaseUri: defaultDatabaseUri,
	}
}

func createTracesExporter(
	ctx context.Context,
	settings exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	config := cfg.(*Config)
	exp, err := newTracesExporter(*config, settings)
	if err != nil {
		return nil, fmt.Errorf("cannot create TraceExporter: %s", err)
	}
	return exporterhelper.NewTracesExporter(
		ctx,
		settings,
		config,
		exp.tracesPusher,
		exporterhelper.WithRetry(config.RetrySettings),
		exporterhelper.WithTimeout(config.TimeoutSettings),
	)
}
