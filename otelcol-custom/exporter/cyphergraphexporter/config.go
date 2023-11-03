package cyphergraphexporter // import "github.com/clement-casse/Playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"go.opentelemetry.io/collector/config/configopaque"
)

type Config struct {
	Username    string              `mapstructure:"username,omitempty"`
	Password    configopaque.String `mapstructure:"password,omitempty"`
	DatabaseUri string              `mapstructure:"db_uri,omitempty"`
}
