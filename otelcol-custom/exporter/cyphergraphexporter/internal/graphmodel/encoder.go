package graphmodel

// Encoder defines the configuration required to encode traces in a hierarchical property graph
type Encoder struct {
	// resourceMap associates the OpenTelemetry Attributes Keys that uniquely identify resources of the given label.
	resourceMap map[string]string
}

// NewEncoder creates a graph encoder with the provided parameters
func NewEncoder(labelFromAttr map[string]string) *Encoder {
	return &Encoder{
		resourceMap: labelFromAttr,
	}
}
