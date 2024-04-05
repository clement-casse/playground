package graphmodel

import (
	"context"
	_ "embed"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

//go:embed mergeResources.cypher
var cypherQueryNewResource string

func (e *Encoder) MergeResource(ctx context.Context, tx neo4j.ManagedTransaction, r *pcommon.Resource) error {
	// Forge a list of map[string]any to be sent as query parameter of neo4j query
	// each entry in rParams is a distinct resource that will be MERGEd with existing
	// one in the graph. Merging occurs on resources sharing both the same "label" and "id" keys.
	// Each resource is a map[string]any expecting the following keys:
	// - "label" designates the kind of resource e.g. a K8S pod, a VM, a container
	// - "id" designates its unique identifier in the system
	// - "props" is another map[string]string with additional parameter to be added to the resource as additional context.
	// - "contained_in" links to the node that contains this resource by identifying it by both its label and id
	rParams := make([]map[string]any, 0, len(e.resourceMap))
	for attrKey, re := range e.resourceMap {
		if attrValue, ok := r.Attributes().Get(string(attrKey)); ok {
			props := make(map[string]any, len(re.AdditionalProps))
			for _, ap := range re.AdditionalProps {
				if value, ok := r.Attributes().Get(string(ap)); ok {
					props[string(ap)] = value.AsString()
				}
			}
			var containedIn []map[string]any
			if rlist, ok := e.containmentOrder[re.ResourceType]; ok {
				containedIn = make([]map[string]any, 0, len(rlist))
				for _, r := range rlist {
					// Reverse lookup of e.resourceMap to find which id has the given labels
					// Should be cached in the encoder to avoid looping over resources again
					// in an inner for loop
					for key, rMapper := range e.resourceMap {
						if rMapper.ResourceType == r {
							containedIn = append(containedIn, map[string]any{
								"label": r,
								"id":    string(key),
							})
							break
						}
					}
				}
			} else {
				containedIn = []map[string]any{}
			}
			rParams = append(rParams, map[string]any{
				"label":        re.ResourceType,
				"id":           attrValue.AsString(),
				"props":        props,
				"contained_in": containedIn,
			})
		}
	}
	// Execute a cypher query on the list of resources that UNWINDs the full list.
	result, err := tx.Run(ctx, cypherQueryNewResource, map[string]any{"resources": rParams})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}
