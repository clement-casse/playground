package cyphergraphexporter

import (
	"fmt"

	neo4jconfig "github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
	"go.opentelemetry.io/collector/config/configretry"
	"go.uber.org/zap"
)

func withBackOffConfig(boc *configretry.BackOffConfig) func(*neo4jconfig.Config) {
	return func(cfg *neo4jconfig.Config) {
		if boc == nil {
			return
		}
		if !boc.Enabled {
			return
		}
		cfg.MaxTransactionRetryTime = boc.MaxInterval
	}
}

func withUserAgent(ua string) func(*neo4jconfig.Config) {
	return func(cfg *neo4jconfig.Config) {
		cfg.UserAgent = ua
	}
}

func withLogger(zl *zap.Logger) func(*neo4jconfig.Config) {
	l := &logWrapper{zsl: zl.Sugar()}
	return func(cfg *neo4jconfig.Config) {
		cfg.Log = l
	}
}

// logWrapper wraps the zap logger to comply neo4j Log interface "github.com/neo4j/neo4j-go-driver/v5/neo4j/log"
type logWrapper struct {
	zsl *zap.SugaredLogger
}

func (l logWrapper) Error(name, id string, err error) {
	l.zsl.Named(name).Errorf("id=%s %s", id, err)
}

func (l logWrapper) Warnf(name, id, msg string, args ...any) {
	l.zsl.Named(name).Warnf(fmt.Sprintf("id=%s %s", id, msg), args...)
}

func (l logWrapper) Infof(name, id, msg string, args ...any) {
	l.zsl.Named(name).Infof(fmt.Sprintf("id=%s %s", id, msg), args...)
}

func (l logWrapper) Debugf(name, id, msg string, args ...any) {
	l.zsl.Named(name).Debugf(fmt.Sprintf("id=%s %s", id, msg), args...)
}
