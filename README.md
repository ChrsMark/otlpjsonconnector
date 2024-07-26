### otlpjsonconnector

*Note*: This component has been upstreamed. Can be found at https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/connector/otlpjsonconnector.

This connector receives log records from a pipeline which contain `otlpjson` in their Body fields
and parse the `otlpjson` extracting Logs, Metrics and Traces according to the content.

#### How to use it

1. Build a custom Collector with the component included (set the path properly in `replace` section): `builder --config=builder-config.yaml`
2. Start the Collector: `./otelcol-dev/otelcol-dev-bin --config=config.yaml`
3. Write a sample log on the target file:
```console
echo '{"resourceLogs":[{"resource":{"attributes":[{"key":"resource-attr","value":{"stringValue":"resource-attr-val-1"}}]},"scopeLogs":[{"scope":{},"logRecords":[{"timeUnixNano":"1581452773000000789","severityNumber":9,"severityText":"Info","body":{"stringValue":"This is a log message"},"attributes":[{"key":"app","value":{"stringValue":"server"}},{"key":"instance_num","value":{"intValue":"1"}}],"droppedAttributesCount":1,"traceId":"08040201000000000000000000000000","spanId":"0102040800000000"},{"timeUnixNano":"1581452773000000789","severityNumber":9,"severityText":"Info","body":{"stringValue":"something happened"},"attributes":[{"key":"customer","value":{"stringValue":"acme"}},{"key":"env","value":{"stringValue":"dev"}}],"droppedAttributesCount":1,"traceId":"","spanId":""}]}]}]}' >> /var/log/pods/prod_my-target-pod_49cc7c1fd3702c40b2686ea7486091d3/my-target-pod/1.log
```
4. Verify the console:
```console
2024-07-19T10:06:09.850+0300	info	ResourceLog #0
Resource SchemaURL: 
Resource attributes:
     -> resource-attr: Str(resource-attr-val-1)
ScopeLogs #0
ScopeLogs SchemaURL: 
InstrumentationScope  
LogRecord #0
ObservedTimestamp: 1970-01-01 00:00:00 +0000 UTC
Timestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
SeverityText: Info
SeverityNumber: Info(9)
Body: Str(This is a log message)
Attributes:
     -> app: Str(server)
     -> instance_num: Int(1)
Trace ID: 08040201000000000000000000000000
Span ID: 0102040800000000
Flags: 0
LogRecord #1
ObservedTimestamp: 1970-01-01 00:00:00 +0000 UTC
Timestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
SeverityText: Info
SeverityNumber: Info(9)
Body: Str(something happened)
Attributes:
     -> customer: Str(acme)
     -> env: Str(dev)
Trace ID: 
Span ID: 
Flags: 0
	{"kind": "exporter", "data_type": "logs", "name": "debug"}
```

Apply the same for metrics and traces:

- extract traces:
```console
echo '{"resourceSpans":[{"resource":{"attributes":[{"key":"resource-attr","value":{"stringValue":"resource-attr-val-1"}}]},"scopeSpans":[{"scope":{},"spans":[{"traceId":"","spanId":"","parentSpanId":"","name":"operationA","startTimeUnixNano":"1581452772000000321","endTimeUnixNano":"1581452773000000789","droppedAttributesCount":1,"events":[{"timeUnixNano":"1581452773000000123","name":"event-with-attr","attributes":[{"key":"span-event-attr","value":{"stringValue":"span-event-attr-val"}}],"droppedAttributesCount":2},{"timeUnixNano":"1581452773000000123","name":"event","droppedAttributesCount":2}],"droppedEventsCount":1,"status":{"message":"status-cancelled","code":2}},{"traceId":"","spanId":"","parentSpanId":"","name":"operationB","startTimeUnixNano":"1581452772000000321","endTimeUnixNano":"1581452773000000789","links":[{"traceId":"","spanId":"","attributes":[{"key":"span-link-attr","value":{"stringValue":"span-link-attr-val"}}],"droppedAttributesCount":4},{"traceId":"","spanId":"","droppedAttributesCount":1}],"droppedLinksCount":3,"status":{}}]}]}]}' >> /var/log/pods/prod_my-target-pod_49cc7c1fd3702c40b2686ea7486091d3/my-target-pod/1.log
```

output:
```console
2024-07-19T10:41:23.592+0300	info	ResourceSpans #0
Resource SchemaURL: 
Resource attributes:
     -> resource-attr: Str(resource-attr-val-1)
ScopeSpans #0
ScopeSpans SchemaURL: 
InstrumentationScope  
Span #0
    Trace ID       : 
    Parent ID      : 
    ID             : 
    Name           : operationA
    Kind           : Unspecified
    Start time     : 2020-02-11 20:26:12.000000321 +0000 UTC
    End time       : 2020-02-11 20:26:13.000000789 +0000 UTC
    Status code    : Error
    Status message : status-cancelled
Events:
SpanEvent #0
     -> Name: event-with-attr
     -> Timestamp: 2020-02-11 20:26:13.000000123 +0000 UTC
     -> DroppedAttributesCount: 2
     -> Attributes::
          -> span-event-attr: Str(span-event-attr-val)
SpanEvent #1
     -> Name: event
     -> Timestamp: 2020-02-11 20:26:13.000000123 +0000 UTC
     -> DroppedAttributesCount: 2
Span #1
    Trace ID       : 
    Parent ID      : 
    ID             : 
    Name           : operationB
    Kind           : Unspecified
    Start time     : 2020-02-11 20:26:12.000000321 +0000 UTC
    End time       : 2020-02-11 20:26:13.000000789 +0000 UTC
    Status code    : Unset
    Status message : 
Links:
SpanLink #0
     -> Trace ID: 
     -> ID: 
     -> TraceState: 
     -> DroppedAttributesCount: 4
     -> Attributes::
          -> span-link-attr: Str(span-link-attr-val)
SpanLink #1
     -> Trace ID: 
     -> ID: 
     -> TraceState: 
     -> DroppedAttributesCount: 1
```

- extracts metrics:
```console
echo '{"resourceMetrics":[{"resource":{"attributes":[{"key":"resource-attr","value":{"stringValue":"resource-attr-val-1"}}]},"scopeMetrics":[{"scope":{},"metrics":[{"name":"counter-int","unit":"1","sum":{"dataPoints":[{"attributes":[{"key":"label-1","value":{"stringValue":"label-value-1"}}],"startTimeUnixNano":"1581452773000000789","timeUnixNano":"1581452773000000789","asInt":"123"},{"attributes":[{"key":"label-2","value":{"stringValue":"label-value-2"}}],"startTimeUnixNano":"1581452772000000321","timeUnixNano":"1581452773000000789","asInt":"456"}],"aggregationTemporality":2,"isMonotonic":true}},{"name":"counter-int","unit":"1","sum":{"dataPoints":[{"attributes":[{"key":"label-1","value":{"stringValue":"label-value-1"}}],"startTimeUnixNano":"1581452772000000321","timeUnixNano":"1581452773000000789","asInt":"123"},{"attributes":[{"key":"label-2","value":{"stringValue":"label-value-2"}}],"startTimeUnixNano":"1581452772000000321","timeUnixNano":"1581452773000000789","asInt":"456"}],"aggregationTemporality":2,"isMonotonic":true}}]}]}]}' >> /var/log/pods/prod_my-target-pod_49cc7c1fd3702c40b2686ea7486091d3/my-target-pod/1.log 
```

output:
```console
2024-07-19T10:40:42.791+0300	info	ResourceMetrics #0
Resource SchemaURL: 
Resource attributes:
     -> resource-attr: Str(resource-attr-val-1)
ScopeMetrics #0
ScopeMetrics SchemaURL: 
InstrumentationScope  
Metric #0
Descriptor:
     -> Name: counter-int
     -> Description: 
     -> Unit: 1
     -> DataType: Sum
     -> IsMonotonic: true
     -> AggregationTemporality: Cumulative
NumberDataPoints #0
Data point attributes:
     -> label-1: Str(label-value-1)
StartTimestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
Timestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
Value: 123
NumberDataPoints #1
Data point attributes:
     -> label-2: Str(label-value-2)
StartTimestamp: 2020-02-11 20:26:12.000000321 +0000 UTC
Timestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
Value: 456
Metric #1
Descriptor:
     -> Name: counter-int
     -> Description: 
     -> Unit: 1
     -> DataType: Sum
     -> IsMonotonic: true
     -> AggregationTemporality: Cumulative
NumberDataPoints #0
Data point attributes:
     -> label-1: Str(label-value-1)
StartTimestamp: 2020-02-11 20:26:12.000000321 +0000 UTC
Timestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
Value: 123
NumberDataPoints #1
Data point attributes:
     -> label-2: Str(label-value-2)
StartTimestamp: 2020-02-11 20:26:12.000000321 +0000 UTC
Timestamp: 2020-02-11 20:26:13.000000789 +0000 UTC
Value: 456
```