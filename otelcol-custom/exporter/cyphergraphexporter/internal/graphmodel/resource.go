package graphmodel

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

const (
	// language=cypher
	cypherQueryNewResource = `MERGE (r:Resource {id: $id, type: $type})`
)

func (e *Encoder) MergeResource(ctx context.Context, tx neo4j.ManagedTransaction, r *pcommon.Resource) error {
	for attrKey, label := range e.resourceMap {
		if attrValue, ok := r.Attributes().AsRaw()[attrKey]; ok {
			result, err := tx.Run(ctx, cypherQueryNewResource, map[string]any{
				"id":   attrValue,
				"type": label,
			})
			if err != nil {
				return err
			}
			_, err = result.Consume(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
