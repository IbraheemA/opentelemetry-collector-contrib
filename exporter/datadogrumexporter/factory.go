// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package datadogrumexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter"

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogrumexporter/internal/metadata"
	"runtime"
	"sync"
	"time"

	pb "github.com/DataDog/datadog-agent/pkg/proto/pbgo/trace"
	"github.com/DataDog/datadog-agent/pkg/trace/writer"
	"github.com/DataDog/opentelemetry-mapping-go/pkg/inframetadata"
	"github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes"
	"github.com/DataDog/opentelemetry-mapping-go/pkg/otlp/attributes/source"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/featuregate"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const metadataReporterPeriod = 30 * time.Minute

type factory struct {
	onceMetadata sync.Once

	onceProvider   sync.Once
	sourceProvider source.Provider
	providerErr    error

	onceReporter     sync.Once
	onceStopReporter sync.Once
	reporter         *inframetadata.Reporter
	reporterErr      error

	onceAttributesTranslator sync.Once
	attributesTranslator     *attributes.Translator
	attributesErr            error

	registry *featuregate.Registry
}

func (f *factory) AttributesTranslator(set component.TelemetrySettings) (*attributes.Translator, error) {
	f.onceAttributesTranslator.Do(func() {
		f.attributesTranslator, f.attributesErr = attributes.NewTranslator(set)
	})
	return f.attributesTranslator, f.attributesErr
}

func newFactoryWithRegistry(registry *featuregate.Registry) exporter.Factory {
	f := &factory{registry: registry}
	return exporter.NewFactory(
		metadata.Type,
		f.createDefaultConfig,
		exporter.WithTraces(f.createTracesExporter, metadata.TracesStability))
}

// NewFactory creates a Datadog exporter factory
func NewFactory() exporter.Factory {
	return newFactoryWithRegistry(featuregate.GlobalRegistry())
}

func defaultClientConfig() confighttp.ClientConfig {
	// do not use NewDefaultClientConfig for backwards-compatibility
	return confighttp.ClientConfig{
		Timeout: 15 * time.Second,
	}
}

// createDefaultConfig creates the default exporter configuration
func (f *factory) createDefaultConfig() component.Config {
	return &Config{}
}

// checkAndCastConfig checks the configuration type and its warnings, and casts it to
// the Datadog Config struct.
func checkAndCastConfig(c component.Config, logger *zap.Logger) *Config {
	cfg, ok := c.(*Config)
	if !ok {
		panic("programming error: config structure is not of type *datadogrumexporter.Config")
	}
	//cfg.logWarnings(logger)
	return cfg
}

func (f *factory) consumeStatsPayload(ctx context.Context, wg *sync.WaitGroup, statsIn <-chan []byte, statsWriter *writer.DatadogStatsWriter, tracerVersion string, agentVersion string, logger *zap.Logger) {
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-statsIn:
					sp := &pb.StatsPayload{}

					err := proto.Unmarshal(msg, sp)
					if err != nil {
						logger.Error("failed to unmarshal stats payload", zap.Error(err))
						continue
					}
					for _, csp := range sp.Stats {
						if csp.TracerVersion == "" {
							csp.TracerVersion = tracerVersion
						}
					}
					// The DD Connector doesn't set the agent version, so we'll set it here
					sp.AgentVersion = agentVersion
					statsWriter.Write(sp)
				}
			}
		}()
	}
}

// createTracesExporter creates a trace exporter based on this config.
func (f *factory) createTracesExporter(
	ctx context.Context,
	set exporter.Settings,
	c component.Config,
) (exporter.Traces, error) {
	cfg := checkAndCastConfig(c, set.TelemetrySettings.Logger)

	var (
		pusher consumer.ConsumeTracesFunc
		stop   component.ShutdownFunc
		wg     sync.WaitGroup // waits for agent to exit
	)

	ctx, cancel := context.WithCancel(ctx)
	// cancel() runs on shutdown

	tracex, err2 := newTracesExporter(ctx, set, cfg)
	if err2 != nil {
		cancel()
		wg.Wait() // then wait for shutdown
		return nil, err2
	}
	//pusher = tracex.consumeTraces
	stop = func(context.Context) error {
		cancel() // first cancel context
		return nil
	}

	pusher = tracex.consumeTraces

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		pusher,
		// explicitly disable since we rely on http.Client timeout logic.
		//exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0 * time.Second}),
		// We don't do retries on traces because of deduping concerns on APM Events.
		exporterhelper.WithRetry(configretry.BackOffConfig{Enabled: false}),
		exporterhelper.WithShutdown(stop),
	)
}
