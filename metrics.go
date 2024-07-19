package otlpjsonconnector

import (
	"context"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// schema for connector
type connectorMetrics struct {
	config          Config
	metricsConsumer consumer.Metrics
	logger          *zap.Logger

	component.StartFunc
	component.ShutdownFunc
}

// newMetricsConnector is a function to create a new connector for metrics extraction
func newMetricsConnector(ctx context.Context, logger *zap.Logger, config component.Config, metricsConsumer consumer.Metrics) (*connectorMetrics, error) {
	logger.Info("Building otlpjson connector for metrics")
	cfg := config.(*Config)

	return &connectorMetrics{
		config:          *cfg,
		logger:          logger,
		metricsConsumer: metricsConsumer,
	}, nil
}

// Capabilities implements the consumer interface.
func (c *connectorMetrics) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// ConsumeLogs method is called for each instance of a log sent to the connector
func (c *connectorMetrics) ConsumeLogs(ctx context.Context, td plog.Logs) error {
	// loop through the levels of logs
	metricsUnmarshaler := &pmetric.JSONUnmarshaler{}
	for i := 0; i < td.ResourceLogs().Len(); i++ {
		li := td.ResourceLogs().At(i)
		for j := 0; j < li.ScopeLogs().Len(); j++ {
			logRecord := li.ScopeLogs().At(j)
			for k := 0; k < logRecord.LogRecords().Len(); k++ {
				lRecord := logRecord.LogRecords().At(k)
				token := lRecord.Body()
				var m pmetric.Metrics
				m, _ = metricsUnmarshaler.UnmarshalMetrics([]byte(token.AsString()))
				return c.metricsConsumer.ConsumeMetrics(ctx, m)
			}
		}
	}
	return nil
}
