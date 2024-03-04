package main

import (
	"context"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/polarsignals/otel-profiling-go/otelgrpcprofiling"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/reflection"
)

func main() {
	otel.SetTracerProvider(initTracer())

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen; %v", err)
	}

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

	echo.RegisterEchoServer(
		server,
		&EchoServer{},
	)

	reflection.Register(server)

	server.Serve(lis)
}

type EchoServer struct {
	echo.UnimplementedEchoServer
}

func (e *EchoServer) UnaryEcho(ctx context.Context, r *echo.EchoRequest) (*echo.EchoResponse, error) {
	fib := fibonacci(42)
	return &echo.EchoResponse{
		Message: fmt.Sprintf("UnaryEcho (fib(42) == %d): %s", fib, r.Message),
	}, nil
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func initTracer() *trace.TracerProvider {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp
}
