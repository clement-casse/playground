package cyphergraphexporter

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	for _, tt := range []struct {
		id                   component.ID
		expected             component.Config
		configValidateAssert assert.ErrorAssertionFunc
	}{
		{
			id:                   component.NewIDWithName(metadata.Type, ""),
			expected:             createDefaultConfig(),
			configValidateAssert: assert.NoError,
		}, {
			id: component.NewIDWithName(metadata.Type, "withnoauth"),
			expected: &Config{
				DatabaseURI: "bolt://neo4j-host:7687",
				UserAgent:   defaultUserAgent,
			},
			configValidateAssert: assert.NoError,
		}, {
			id: component.NewIDWithName(metadata.Type, "withbasicauth"),
			expected: &Config{
				DatabaseURI: "bolt://neo4j-host:7687",
				Username:    "neo4j",
				Password:    configopaque.String("password"),
				UserAgent:   defaultUserAgent,
			},
			configValidateAssert: assert.NoError,
		}, {
			id: component.NewIDWithName(metadata.Type, "withbearertoken"),
			expected: &Config{
				DatabaseURI: "bolt://neo4j-host:7687",
				BearerToken: configopaque.String("somevalue"),
				UserAgent:   defaultUserAgent,
			},
			configValidateAssert: assert.NoError,
		}, {
			id: component.NewIDWithName(metadata.Type, "withkerberosticket"),
			expected: &Config{
				DatabaseURI:    "bolt://neo4j-host:7687",
				KerberosTicket: configopaque.String("somevalue"),
				UserAgent:      defaultUserAgent,
			},
			configValidateAssert: assert.NoError,
		}, {
			id: component.NewIDWithName(metadata.Type, "withcustomua"),
			expected: &Config{
				DatabaseURI: defaultDatabaseURI,
				UserAgent:   "testUserAgent",
			},
			configValidateAssert: assert.NoError,
		}, {
			id: component.NewIDWithName(metadata.Type, "withbasicandbearer"),
			expected: &Config{
				DatabaseURI: defaultDatabaseURI,
				UserAgent:   defaultUserAgent,
				Username:    "neo4j",
				Password:    configopaque.String("password"),
				BearerToken: configopaque.String("somevalue"),
			},
			configValidateAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, errMultipleAuthMethod.Error())
			},
		}, {
			id: component.NewIDWithName(metadata.Type, "withabearerandkerb"),
			expected: &Config{
				DatabaseURI:    defaultDatabaseURI,
				UserAgent:      defaultUserAgent,
				BearerToken:    configopaque.String("somevalue"),
				KerberosTicket: configopaque.String("somevalue"),
			},
			configValidateAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, errMultipleAuthMethod.Error())
			},
		}, {
			id: component.NewIDWithName(metadata.Type, "withallauthmethods"),
			expected: &Config{
				DatabaseURI:    defaultDatabaseURI,
				UserAgent:      defaultUserAgent,
				Username:       "neo4j",
				Password:       configopaque.String("password"),
				BearerToken:    configopaque.String("somevalue"),
				KerberosTicket: configopaque.String("somevalue"),
			},
			configValidateAssert: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorContains(t, err, errMultipleAuthMethod.Error())
			},
		},
	} {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			vv := component.ValidateConfig(cfg)
			tt.configValidateAssert(t, vv)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
