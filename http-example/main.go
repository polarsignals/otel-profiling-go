package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	otelprof "github.com/polarsignals/otel-profiling-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	otel.SetTracerProvider(initTracer())

	err := http.ListenAndServe(":3000",
		otelhttp.NewHandler(http.HandlerFunc(fibHandler), "fibHandler"),
	)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func fibHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.GetTracerProvider().Tracer("")
	_, span := tracer.Start(r.Context(), "fibHandler")
	defer span.End()

	w.Write([]byte(fmt.Sprintf("fib %d\n", fibonacci(42))))
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
	otel.SetTextMapPropagator(
		otelprof.NewTextMapPropagatorWithProfiling(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			)))
	return tp
}
