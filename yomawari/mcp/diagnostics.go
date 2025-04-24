package mcp

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "Experimental.ModelContextProtocol"
	metricName = "Experimental.ModelContextProtocol"
)

var (
	Tracer                       = otel.Tracer(tracerName)
	meter                        = otel.Meter(metricName)
	ShortSecondsBucketBoundaries = []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10}
	LongSecondsBucketBoundaries  = []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 30, 60, 120, 300}
)

var Propagator = otel.GetTextMapPropagator()

func CreateDurationHistogram(name string, description string, longBuckets bool) metric.Float64Histogram {
	boundaries := ShortSecondsBucketBoundaries
	if longBuckets {
		boundaries = LongSecondsBucketBoundaries
	}
	m, _ := meter.Float64Histogram(
		name,
		metric.WithUnit("s"),
		metric.WithDescription(description),
		metric.WithExplicitBucketBoundaries(boundaries...),
	)
	return m
}

func StartServerSpan(ctx context.Context, name string, carrier propagation.TextMapCarrier) (context.Context, trace.Span) {
	parentCtx := Propagator.Extract(ctx, carrier)
	link := trace.Link{SpanContext: trace.SpanContextFromContext(parentCtx)}
	return Tracer.Start(parentCtx, name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithLinks(link),
	)
}

func StartSpanWithJsonRpcData(ctx context.Context, name string, message transport.IJsonRpcMessage) (context.Context, trace.Span) {
	var carrier propagation.TextMapCarrier = propagation.MapCarrier{}
	switch re := message.(type) {
	case *transport.JsonRpcRequest:
		carrier = re
	case *transport.JsonRpcNotification:
		carrier = re
	}
	parentCtx := Propagator.Extract(ctx, carrier)
	link := trace.Link{SpanContext: trace.SpanContextFromContext(parentCtx)}
	return Tracer.Start(parentCtx, name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithLinks(link),
	)
}

func PropagatorInject(ctx context.Context, message transport.IJsonRpcMessage) {
	var carrier propagation.TextMapCarrier = propagation.MapCarrier{}
	switch re := message.(type) {
	case *transport.JsonRpcRequest:
		carrier = re
	case *transport.JsonRpcNotification:
		carrier = re
	}
	Propagator.Inject(ctx, carrier)
}
