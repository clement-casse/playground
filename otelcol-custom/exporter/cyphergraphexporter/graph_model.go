package cyphergraphexporter

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

const (
	// language=cypher
	cypherQueryNewResource = `MERGE (r:Resource {})`
)

func mergeResource(ctx context.Context, tx neo4j.ManagedTransaction, r *pcommon.Resource) error {
	result, err := tx.Run(ctx, cypherQueryNewResource, r.Attributes().AsRaw())
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}
