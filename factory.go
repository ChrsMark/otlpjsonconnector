package otlpjsonconnector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/consumer"
)

var (
	typeStr = component.MustNewType("otlpjson")
)

// NewFactory returns a ConnectorFactory.
func NewFactory() connector.Factory {
	return connector.NewFactory(
		typeStr, // metadata.Type,
		createDefaultConfig,
		connector.WithLogsToTraces(createTracesConnector, component.StabilityLevelAlpha),
		connector.WithLogsToMetrics(createMetricsConnector, component.StabilityLevelAlpha),
		connector.WithLogsToLogs(createLogsConnector, component.StabilityLevelAlpha),
	)
}

// createDefaultConfig creates the default configuration.
func createDefaultConfig() component.Config {
	return &Config{}
}

// createLogsToTraces defines the consumer type of the connector
// We want to consume logs and export logs, therefore, define nextConsumer as logs, since consumer is the next component in the pipeline
func createLogsConnector(
	ctx context.Context,
	params connector.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (connector.Logs, error) {
	return newLogsConnector(ctx, params.Logger, cfg, nextConsumer)
}

// createLogsToTraces defines the consumer type of the connector
// We want to consume logs and export logs, therefore, define nextConsumer as logs, since consumer is the next component in the pipeline
func createTracesConnector(
	ctx context.Context,
	params connector.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (connector.Logs, error) {
	return newTracesConnector(ctx, params.Logger, cfg, nextConsumer)
}

// createLogsToTraces defines the consumer type of the connector
// We want to consume logs and export logs, therefore, define nextConsumer as logs, since consumer is the next component in the pipeline
func createMetricsConnector(
	ctx context.Context,
	params connector.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (connector.Logs, error) {
	return newMetricsConnector(ctx, params.Logger, cfg, nextConsumer)
}
