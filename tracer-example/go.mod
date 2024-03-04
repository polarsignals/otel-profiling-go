module github.com/polarsignals/otel-profiling-go/example

go 1.21.6

toolchain go1.22.0

require (
	github.com/polarsignals/otel-profiling-go v0.1.0
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.1
	go.opentelemetry.io/otel/sdk v1.21.0
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace github.com/polarsignals/otel-profiling-go => ../
