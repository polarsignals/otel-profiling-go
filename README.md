# otel-profiling-go

This package provides an integration between distributed tracing via OpenTelemetry with Profiling data collected via eBPF by [Parca Agent](https://github.com/parca-dev/parca-agent). The best thing about this is that it isn't actually a deep integration with Parca Agent, it just puts the trace ID into a place that Parca Agent will know how to find.

> Note: Currently Parca Agent only supports reading the trace ID in Go 1.22.

More specifically it provides three ways to do so:

* An OpenTelemetry `Tracer` implementation
* A gRPC middleware
* An HTTP middleware

What these provide are ways to automatically [set Go's goroutine labels](https://pkg.go.dev/runtime/pprof#SetGoroutineLabels). Goroutine labels are then accessible to Parca Agent via eBPF, as they are stored in thread-local-store, which the Go runtime manages, which has a well-known layout, so it is easy for the Parca Agent to know how to read them.

## Using the Tracer

The tracer is the simplest way to get started, as it can be used as a drop-in replacement for your existing tracer, and you're all set! The tracer is also the least efficient, as it causes allocations with every new span that's created.

```go
import (
	otelprof "github.com/polarsignals/otel-profiling-go"
)

main() {
	tp := initTracer()
	otel.SetTracerProvider(otelprof.NewTracerProvider(tp))

	// your application...
}
```

See a full example in [`./tracer-example`](./tracer-example).

## Middleware

Our recommendation is to use a middleware approach, where the trace ID is set once per request, causing only constant allocations, instead of allocations per span.

### HTTP

For HTTP use the handler wrapper. Make sure it is called after an initial trace ID has been set on the context.

```go
import (
	"github.com/polarsignals/otel-profiling-go/otelhttpprofiling"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	handler := otelhttp.NewHandler(otelhttpprofiling.Handler(http.HandlerFunc(fibHandler)), "fibHandler")

	// ... actually serve handler
}
```

See a full example in [`./http-example`](./http-example).

### gRPC

For gRPC use the gRPC middleware. Same as the HTTP middleware, ensure that it is after the otel interceptors to ensure a trace ID is already set on the context.

```go
import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

main() {
	otel.SetTracerProvider(initTracer())

	grpcotelprof := otelgrpcprofiling.NewMiddleware()
	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			otelgrpc.UnaryServerInterceptor(),
			grpcotelprof.GrpcUnaryInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			otelgrpc.StreamServerInterceptor(),
			grpcotelprof.GrpcStreamInterceptor,
		),
	)

	// ... register gRPC services, serve it, etc.
}
```

See a full example in [`./grpc-example`](./grpc-example).
