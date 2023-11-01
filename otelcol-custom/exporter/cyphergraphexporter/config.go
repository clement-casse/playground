package cyphergraphexporter // import "github.com/clement-casse/Playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
}
