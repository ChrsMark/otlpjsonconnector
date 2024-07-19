package otlpjsonconnector

import (
	"context"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
)

// schema for connector
type connectorLogs struct {
	config       Config
	logsConsumer consumer.Logs
	logger       *zap.Logger

	component.StartFunc
	component.ShutdownFunc
}

// newLogsConnector is a function to create a new connector for logs extraction
func newLogsConnector(ctx context.Context, logger *zap.Logger, config component.Config, logsConsumer consumer.Logs) (*connectorLogs, error) {
	logger.Info("Building otlpjson connector for logs")
	cfg := config.(*Config)

	return &connectorLogs{
		config:       *cfg,
		logger:       logger,
		logsConsumer: logsConsumer,
	}, nil
}

// Capabilities implements the consumer interface.
func (c *connectorLogs) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// ConsumeLogs method is called for each instance of a log sent to the connector
func (c *connectorLogs) ConsumeLogs(ctx context.Context, td plog.Logs) error {
	// loop through the levels of logs
	logsUnmarshaler := &plog.JSONUnmarshaler{}
	for i := 0; i < td.ResourceLogs().Len(); i++ {
		li := td.ResourceLogs().At(i)
		for j := 0; j < li.ScopeLogs().Len(); j++ {
			logRecord := li.ScopeLogs().At(j)
			for k := 0; k < logRecord.LogRecords().Len(); k++ {
				lRecord := logRecord.LogRecords().At(k)
				token := lRecord.Body()
				var l plog.Logs
				l, _ = logsUnmarshaler.UnmarshalLogs([]byte(token.AsString()))
				return c.logsConsumer.ConsumeLogs(ctx, l)
			}
		}
	}
	return nil
}
