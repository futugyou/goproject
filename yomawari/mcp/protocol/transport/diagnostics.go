package transport

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "Experimental.ModelContextProtocol"
)

var (
	Tracer = otel.Tracer(tracerName)
)

var Propagator = otel.GetTextMapPropagator()

func StartSpanWithJsonRpcData(ctx context.Context, name string, message IJsonRpcMessage) (context.Context, trace.Span) {
	var carrier propagation.TextMapCarrier = propagation.MapCarrier{}
	switch re := message.(type) {
	case *JsonRpcRequest:
		carrier = re
	case *JsonRpcNotification:
		carrier = re
	}
	parentCtx := Propagator.Extract(ctx, carrier)
	link := trace.Link{SpanContext: trace.SpanContextFromContext(parentCtx)}
	return Tracer.Start(parentCtx, name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithLinks(link),
	)
}

func PropagatorInject(ctx context.Context, message IJsonRpcMessage) {
	var carrier propagation.TextMapCarrier = propagation.MapCarrier{}
	switch re := message.(type) {
	case *JsonRpcRequest:
		carrier = re
	case *JsonRpcNotification:
		carrier = re
	}
	Propagator.Inject(ctx, carrier)
}
