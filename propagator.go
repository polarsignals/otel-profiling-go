package otelprofiling

import (
	"context"
	"runtime/pprof"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type propagator struct {
	propagation.TextMapPropagator
}

func (p propagator) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	ctx := p.TextMapPropagator.Extract(parent, carrier)
	traceID := trace.SpanContextFromContext(ctx).TraceID().String()
	ctx = pprof.WithLabels(ctx, pprof.Labels("otel.traceid", traceID))
	pprof.SetGoroutineLabels(ctx)
	return ctx
}

// NewTextMapPropagatorWithProfiling creates a propagator that annotates pprof
// samples with the otel.traceid label. This allows to establish a relationship
// between pprof profiles and reported tracing spans.
func NewTextMapPropagatorWithProfiling(base propagation.TextMapPropagator) propagation.TextMapPropagator {
	return propagator{TextMapPropagator: base}
}
