package otlpjsonconnector

import (
	"context"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// schema for connector
type connectorTraces struct {
	config         Config
	tracesConsumer consumer.Traces
	logger         *zap.Logger

	component.StartFunc
	component.ShutdownFunc
}

// newTracesConnector is a function to create a new connector for traces extraction
func newTracesConnector(ctx context.Context, logger *zap.Logger, config component.Config, tracesConsumer consumer.Traces) (*connectorTraces, error) {
	logger.Info("Building otlpjson connector for traces")
	cfg := config.(*Config)

	return &connectorTraces{
		config:         *cfg,
		logger:         logger,
		tracesConsumer: tracesConsumer,
	}, nil
}

// Capabilities implements the consumer interface.
func (c *connectorTraces) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// ConsumeLogs method is called for each instance of a log sent to the connector
func (c *connectorTraces) ConsumeLogs(ctx context.Context, td plog.Logs) error {
	// loop through the levels of logs
	tracesUnmarshaler := &ptrace.JSONUnmarshaler{}
	for i := 0; i < td.ResourceLogs().Len(); i++ {
		li := td.ResourceLogs().At(i)
		for j := 0; j < li.ScopeLogs().Len(); j++ {
			logRecord := li.ScopeLogs().At(j)
			for k := 0; k < logRecord.LogRecords().Len(); k++ {
				lRecord := logRecord.LogRecords().At(k)
				token := lRecord.Body()
				var t ptrace.Traces
				t, _ = tracesUnmarshaler.UnmarshalTraces([]byte(token.AsString()))
				return c.tracesConsumer.ConsumeTraces(ctx, t)
			}
		}
	}
	return nil
}
