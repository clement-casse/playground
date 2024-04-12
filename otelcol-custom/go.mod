module github.com/clement-casse/playground/otelcol-custom

go 1.21

require (
	github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter v0.0.0
	github.com/mikefarah/yq/v4 v4.43.1
	go.opentelemetry.io/collector/cmd/builder v0.97.0
	go.opentelemetry.io/collector/cmd/mdatagen v0.97.0
)

require (
	github.com/a8m/envsubst v1.4.2 // indirect
	github.com/alecthomas/participle/v2 v2.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/elliotchance/orderedmap v1.5.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.11.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/parsers/yaml v0.1.0 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/providers/env v0.1.0 // indirect
	github.com/knadh/koanf/providers/file v0.1.0 // indirect
	github.com/knadh/koanf/providers/fs v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/neo4j/neo4j-go-driver/v5 v5.19.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.0 // indirect
	github.com/spf13/cobra v1.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	go.opentelemetry.io/collector v0.97.0 // indirect
	go.opentelemetry.io/collector/component v0.97.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.4.0 // indirect
	go.opentelemetry.io/collector/config/configretry v0.97.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.97.0 // indirect
	go.opentelemetry.io/collector/confmap v0.97.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/fileprovider v0.97.0 // indirect
	go.opentelemetry.io/collector/consumer v0.97.0 // indirect
	go.opentelemetry.io/collector/exporter v0.97.0 // indirect
	go.opentelemetry.io/collector/extension v0.97.0 // indirect
	go.opentelemetry.io/collector/pdata v1.4.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/grpc v1.62.1 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/op/go-logging.v1 v1.0.0-20160211212156-b2cb9fa56473 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Replace references to modules that are in this repository with their relateive paths
// so that we always build with current (latest) version of the source code.

replace github.com/clement-casse/playground/otelcol-custom/exporter/cyphergraphexporter => ./exporter/cyphergraphexporter
