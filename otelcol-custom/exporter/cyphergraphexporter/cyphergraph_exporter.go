package cyphergraphexporter // import "github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter"

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type cyphergraphTraceExporter struct {
	driver        neo4j.DriverWithContext
	labelFromAttr map[string]string

	logger *zap.Logger
}

func newTracesExporter(cfg *Config, set exporter.CreateSettings) (cte *cyphergraphTraceExporter, err error) {
	cte = &cyphergraphTraceExporter{logger: set.Logger}
	if err = cfg.Validate(); err != nil {
		return
	}
	var neo4jAuth neo4j.AuthToken
	if cfg.Username != "" {
		neo4jAuth = neo4j.BasicAuth(cfg.Username, string(cfg.Password), "")
	} else if t := string(cfg.BearerToken); t != "" {
		neo4jAuth = neo4j.BearerAuth(t)
	} else if t := string(cfg.KerberosTicket); t != "" {
		neo4jAuth = neo4j.KerberosAuth(t)
	} else {
		neo4jAuth = neo4j.NoAuth()
	}
	cte.driver, err = neo4j.NewDriverWithContext(
		cfg.DatabaseURI,
		neo4jAuth,
		withLogger(set.Logger),
		withUserAgent(cfg.UserAgent),
		withBackOffConfig(&cfg.BackOffConfig),
	)
	if err != nil {
		return
	}
	cte.labelFromAttr = make(map[string]string, len(cfg.ResourceMappers))
	for label, matcher := range cfg.ResourceMappers {
		cte.labelFromAttr[matcher.LabelID] = label
	}
	return
}

func (e *cyphergraphTraceExporter) Start(ctx context.Context, _ component.Host) error {
	err := e.driver.VerifyConnectivity(ctx)
	if err != nil {
		return err
	}
	// TODO think about initializing database indices here.
	return nil
}

func (e *cyphergraphTraceExporter) Shutdown(ctx context.Context) error {
	if e.driver == nil {
		return nil
	}
	return e.driver.Close(ctx)
}

func (e *cyphergraphTraceExporter) tracesPusher(ctx context.Context, td ptrace.Traces) error {
	start := time.Now()
	session := e.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	rSpans := td.ResourceSpans()
	for i := 0; i < rSpans.Len(); i++ {
		e.logger.Sugar().Infof("span %d, Span: %+v", i, rSpans.At(i))
		// _, err := neo4j.ExecuteWrite(ctx, session, func(tx neo4j.ManagedTransaction) (any, error) {
		// 	resource := rSpans.At(i).Resource()
		// 	if err := graphmodel.MergeResource(ctx, tx, &resource); err != nil {
		// 		return nil, err
		// 	}

		// 	return nil, nil
		// })
		// if err != nil {
		// 	e.logger.Error("resources merge error", zap.Error(err))
		// }
	}

	duration := time.Since(start)
	e.logger.Debug("traces inserted", zap.String("duration", duration.String()))
	return nil
}
