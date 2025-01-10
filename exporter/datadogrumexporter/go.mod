module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter

go 1.22.0

require (
	github.com/DataDog/datadog-agent/comp/otelcol/otlp/testutil v0.57.0-devel.0.20240718200853-81bf3b2e412d
	github.com/DataDog/datadog-agent/pkg/trace v0.57.0-devel.0.20240722160158-ad956a31a730
	github.com/DataDog/datadog-api-client-go/v2 v2.29.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata v0.19.0
	github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes v0.19.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/collector/component v0.107.0
	go.opentelemetry.io/collector/config/confighttp v0.107.0
	go.opentelemetry.io/collector/config/confignet v0.107.0
	go.opentelemetry.io/collector/config/configretry v1.13.0
	go.opentelemetry.io/collector/confmap v0.107.0
	go.opentelemetry.io/collector/consumer v0.107.0
	go.opentelemetry.io/collector/exporter v0.107.0
	go.opentelemetry.io/collector/featuregate v1.13.0
	go.opentelemetry.io/collector/pdata v1.13.0
	go.opentelemetry.io/collector/semconv v0.107.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/DataDog/datadog-agent/comp/core/secrets v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/collector/check/defaults v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/env v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/model v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/config/setup v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/proto v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/executable v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/filesystem v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/hostname/validate v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/log v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/optional v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/pointer v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/scrubber v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/system v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/system/socket v0.56.0 // indirect
	github.com/DataDog/datadog-agent/pkg/util/winutil v0.56.0 // indirect
	github.com/DataDog/sketches-go v1.4.6 // indirect
	github.com/DataDog/viper v1.13.5 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.0.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hectane/go-acl v0.0.0-20190604041725-da78bae5fc95 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/cors v1.11.0 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shoenig/test v1.7.1 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/cobra v1.8.1 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tinylib/msgp v1.1.9 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/collector v0.107.0 // indirect
	go.opentelemetry.io/collector/client v1.13.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.107.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.13.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.13.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.107.0 // indirect
	go.opentelemetry.io/collector/config/configtls v1.13.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.107.0 // indirect
	go.opentelemetry.io/collector/consumer/consumerprofiles v0.107.0 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.107.0 // indirect
	go.opentelemetry.io/collector/extension v0.107.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.107.0 // indirect
	go.opentelemetry.io/collector/internal/globalgates v0.107.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.107.0 // indirect
	go.opentelemetry.io/collector/receiver v0.107.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.50.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
