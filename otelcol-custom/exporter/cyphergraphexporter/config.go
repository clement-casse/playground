package cyphergraphexporter // import "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configretry"
)

// Config defines configuration for the Cypher Graph exporter.
type Config struct {
	configretry.BackOffConfig `mapstructure:"retry_on_failure"`

	Username    string              `mapstructure:"username,omitempty"`
	Password    configopaque.String `mapstructure:"password,omitempty"`
	DatabaseURI string              `mapstructure:"db_uri,omitempty"`

	UserAgent string `mapstructure:"user_agent,omitempty"`
}

func (cfg *Config) Validate() error {
	return nil
}
