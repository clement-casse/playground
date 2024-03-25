package cyphergraphexporter // import "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter/internal/metadata"
)

const (
	defaultUsername    = ""
	defaultPassword    = ""
	defaultDatabaseURI = "bolt://localhost:7687"
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
		DatabaseURI: defaultDatabaseURI,
	}
}

func createTracesExporter(
	ctx context.Context,
	settings exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	config := cfg.(*Config)
	exp, err := newTracesExporter(config, settings)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewTracesExporter(ctx, settings, config,
		exp.tracesPusher,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(exp.Start),
		exporterhelper.WithShutdown(exp.Shutdown),
	)
}
