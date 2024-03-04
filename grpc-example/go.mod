module github.com/polarsignals/otel-profiling-go/grpc-example

go 1.21.6

toolchain go1.22.0

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/polarsignals/otel-profiling-go v0.1.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.1
	go.opentelemetry.io/otel/sdk v1.21.0
	google.golang.org/grpc v1.62.0
	google.golang.org/grpc/examples v0.0.0-20240228223710-2a617ca67a6b
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace github.com/polarsignals/otel-profiling-go => ../
