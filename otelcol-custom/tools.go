//go:build tools
// +build tools

package otelcolcustom

import (
	_ "go.opentelemetry.io/collector/cmd/builder"
	_ "go.opentelemetry.io/collector/cmd/mdatagen"

	_ "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"
)
