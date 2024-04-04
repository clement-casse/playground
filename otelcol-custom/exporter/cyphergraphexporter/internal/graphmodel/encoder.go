package graphmodel

import "go.opentelemetry.io/otel/attribute"

// Encoder defines the configuration required to encode traces in a hierarchical property graph
type Encoder struct {
	// resourceMap associates the OpenTelemetry Attributes Keys that uniquely identify resources of the given label.
	resourceMap map[attribute.Key]string
}

// NewEncoder creates a graph encoder with the provided parameters
func NewEncoder(labelFromAttr map[attribute.Key]string) *Encoder {
	return &Encoder{
		resourceMap: labelFromAttr,
	}
}
