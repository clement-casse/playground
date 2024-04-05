package graphmodel

import "go.opentelemetry.io/otel/attribute"

// Encoder defines the configuration required to encode traces in a hierarchical property graph
type Encoder struct {
	// resourceMap associates the OpenTelemetry Attributes Keys that uniquely identify resources of the given label.
	resourceMap map[attribute.Key]ResourceEncoder

	// containmentOrder represents the graph of resource containment as an adjacency list:
	// e.g.
	//    containmentOrder := map[string][]string{
	//        "k8s.pod":                 {"k8s.node"},
	//        "k8s.node":                {"k8s.cluster", "cloud.availability.zone"},
	//        "cloud.availability.zone": {"cloud.region"},
	//        "cloud.region":            {},
	//        "k8s.cluster":             {},
	//    }
	containmentOrder map[string][]string
}

type ResourceEncoder struct {
	ResourceType    string
	AdditionalProps []attribute.Key
}

// NewEncoder creates a graph encoder with the provided parameters
func NewEncoder(labelFromAttr map[attribute.Key]ResourceEncoder, containmentOrder map[string][]string) *Encoder {
	return &Encoder{
		resourceMap:      labelFromAttr,
		containmentOrder: containmentOrder,
	}
}
