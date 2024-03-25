//go:generate go run go.opentelemetry.io/collector/cmd/builder --config=./builder-config.yaml

// Package otelcolcustom defines a demonstration OpenTelemetry Collector built with the curent
// versions of the custom modules that are part of this project.
//
// The collector is build with OpenTelemetry Collector Builder (ocb), its components have been
// defined in the manifest ./builder-config.yaml.
package otelcolcustom
