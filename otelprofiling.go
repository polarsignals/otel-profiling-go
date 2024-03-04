package otelprofiling

import (
	"context"
	"runtime/pprof"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// tracerProvider satisfies open telemetry TracerProvider interface.
type tracerProvider struct {
	noop.TracerProvider
	tp trace.TracerProvider
}

type Option func(*tracerProvider)

// NewTracerProvider creates a new tracer provider that annotates pprof
// samples with span_id label. This allows to establish a relationship
// between pprof profiles and reported tracing spans.
func NewTracerProvider(tp trace.TracerProvider, options ...Option) trace.TracerProvider {
	p := tracerProvider{
		tp: tp,
	}
	for _, o := range options {
		o(&p)
	}
	return &p
}

func (w *tracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return &profileTracer{p: w, tr: w.tp.Tracer(name, opts...)}
}

type profileTracer struct {
	noop.Tracer
	p  *tracerProvider
	tr trace.Tracer
}

func (w *profileTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := w.tr.Start(ctx, spanName, opts...)
	spanCtx := span.SpanContext()
	traceID := spanCtx.TraceID().String()

	s := spanWrapper{
		Span: span,
		ctx:  ctx,
		p:    w.p,
	}

	if traceID != "" {
		ctx = pprof.WithLabels(ctx, pprof.Labels("otel.traceid", traceID))
		pprof.SetGoroutineLabels(ctx)
	}
	return ctx, &s
}

type spanWrapper struct {
	trace.Span
	ctx context.Context
	p   *tracerProvider
}

func (s spanWrapper) End(options ...trace.SpanEndOption) {
	s.Span.End(options...)
	pprof.SetGoroutineLabels(s.ctx)
}
