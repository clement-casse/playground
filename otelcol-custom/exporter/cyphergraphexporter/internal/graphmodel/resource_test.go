package graphmodel

import (
	"context"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const databaseEngine = "memgraph"

// resourceFrom returns a [go.opentelemetry.io/collector/pdata/pcommon.Resource] with the provided Key Value pairs.
func resourceFrom(attr []attribute.KeyValue) pcommon.Resource {
	res := pcommon.NewResource()
	for _, kv := range attr {
		switch kv.Value.Type() {
		case attribute.BOOL:
			res.Attributes().PutBool(string(kv.Key), kv.Value.AsBool())
		case attribute.INT64:
			res.Attributes().PutInt(string(kv.Key), kv.Value.AsInt64())
		case attribute.FLOAT64:
			res.Attributes().PutDouble(string(kv.Key), kv.Value.AsFloat64())
		case attribute.STRING:
			res.Attributes().PutStr(string(kv.Key), kv.Value.AsString())
		case attribute.BOOLSLICE, attribute.INT64SLICE, attribute.FLOAT64SLICE, attribute.STRINGSLICE:
			panic("slice attributes are not managed by the resourceFrom function")
		default:
			panic("Type not handled by resourceFrom function: " + kv.Value.Type().String())
		}
	}
	return res
}

var (
	// emptyResource  = pcommon.NewResource()
	// simpleResource = resourceFrom(map[string]string{
	// 	string(semconv.ServiceNameKey): "name-of-the-service",
	// })
	k8sResource = resourceFrom([]attribute.KeyValue{
		semconv.ServiceName("my-deployment"),
		semconv.K8SClusterName("my-cluster"),
		semconv.K8SNamespaceName("my-namespace"),
		semconv.K8SDeploymentName("my-deployment"),
		semconv.K8SReplicaSetName("my-deployment-66cf4d99b5"),
		semconv.K8SPodName("my-deployment-66cf4d99b5-kpqg"),
		semconv.K8SPodUID("7293ca81-d35e-459d-b15a-a8197fbc03e4"),
	})
)

func TestMergeResource(t *testing.T) {
	ctx := context.Background()
	driver, closeFunc, err := InitCypherDBTestContainer(ctx, databaseEngine)
	assert.NoError(t, err)
	defer closeFunc(ctx)

	encoder := NewEncoder(map[attribute.Key]string{
		semconv.K8SPodNameKey:        "k8s.pod",
		semconv.K8SDeploymentNameKey: "k8s.deployment",
		semconv.K8SClusterNameKey:    "k8s.cluster",
	})

	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, encoder.MergeResource(ctx, tx, &k8sResource)
	})
	assert.NoError(t, err)

	_, err = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err2 := tx.Run(ctx, `MATCH (r:Resource) RETURN r.id, r.type`, map[string]any{})
		assert.NoError(t, err2)
		records, err2 := result.Collect(ctx)
		assert.NoError(t, err2)

		for _, rec := range records {
			t.Logf("%+v", rec.Values)
		}
		return nil, nil
	})
	assert.NoError(t, err)
}
