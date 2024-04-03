package graphmodel

import (
	"context"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func resourceFrom(attr map[string]string) pcommon.Resource {
	res := pcommon.NewResource()
	for k, v := range attr {
		res.Attributes().PutStr(k, v)
	}
	return res
}

var (
	databaseEngine = "memgraph"

	// emptyResource  = pcommon.NewResource()
	// simpleResource = resourceFrom(map[string]string{
	// 	string(semconv.ServiceNameKey): "name-of-the-service",
	// })
	k8sResource = resourceFrom(map[string]string{
		string(semconv.ServiceNameKey):       "my-deployment",
		string(semconv.K8SClusterNameKey):    "my-cluster",
		string(semconv.K8SNamespaceNameKey):  "my-namespace",
		string(semconv.K8SDeploymentNameKey): "my-deployment",
		string(semconv.K8SReplicaSetNameKey): "my-deployment-66cf4d99b5",
		string(semconv.K8SPodNameKey):        "my-deployment-66cf4d99b5-kpqg",
		string(semconv.K8SPodUIDKey):         "7293ca81-d35e-459d-b15a-a8197fbc03e4",
	})
)

func TestMergeResource(t *testing.T) {
	ctx := context.Background()
	driver, closeFunc, err := InitCypherDBTestContainer(ctx, databaseEngine)
	assert.NoError(t, err)
	defer closeFunc(ctx)

	encoder := NewEncoder(map[string]string{
		string(semconv.K8SPodNameKey):        "k8s.pod",
		string(semconv.K8SDeploymentNameKey): "k8s.deployment",
		string(semconv.K8SClusterNameKey):    "k8s.cluster",
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
