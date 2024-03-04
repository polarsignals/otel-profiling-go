package otelgrpcprofiling

import (
	"context"
	"runtime/pprof"
	"strings"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Middleware struct {
	includePrefixes []string
	ignorePrefixes  []string
}

type Option func(*Middleware)

func WithPrefix(prefix string) Option {
	return func(m *Middleware) {
		m.includePrefixes = append(m.includePrefixes, prefix)
	}
}

func WithIgnorePrefix(prefix string) Option {
	return func(m *Middleware) {
		m.ignorePrefixes = append(m.ignorePrefixes, prefix)
	}
}

func NewMiddleware(opts ...Option) *Middleware {
	m := &Middleware{}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Middleware) GrpcUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if m.match(info.FullMethod) {
		var (
			traceID = trace.SpanContextFromContext(ctx).TraceID().String()

			res interface{}
			err error
		)

		if traceID != "" {
			pprof.Do(ctx, pprof.Labels("otel.traceid", traceID), func(ctx context.Context) {
				res, err = handler(ctx, req)
			})

			return res, err
		}
	}

	return handler(ctx, req)
}

func (m *Middleware) match(fullMethod string) bool {
	if !m.matchPrefixes(fullMethod) {
		return false
	}

	return !m.matchIgnorePrefixes(fullMethod)
}

func (m *Middleware) matchPrefixes(fullMethod string) bool {
	if len(m.includePrefixes) == 0 {
		return true
	}

	for _, prefix := range m.includePrefixes {
		if strings.HasPrefix(fullMethod, prefix) {
			return true
		}
	}
	return false
}

func (m *Middleware) matchIgnorePrefixes(fullMethod string) bool {
	for _, prefix := range m.ignorePrefixes {
		if strings.HasPrefix(fullMethod, prefix) {
			return true
		}
	}
	return false
}

func (m *Middleware) GrpcStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if m.match(info.FullMethod) {
		var (
			ctx     = ss.Context()
			traceID = trace.SpanContextFromContext(ctx).TraceID().String()

			err error
		)
		pprof.Do(ctx, pprof.Labels("otel.traceid", traceID), func(ctx context.Context) {
			err = handler(srv, ss)
		})

		return err
	}

	return handler(srv, ss)
}
