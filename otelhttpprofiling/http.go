package otelhttpprofiling

import (
	"context"
	"net/http"
	"runtime/pprof"

	"go.opentelemetry.io/otel/trace"
)

func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.SpanContextFromContext(ctx).TraceID().String()

		if traceID != "" {
			pprof.Do(ctx, pprof.Labels("otel.traceid", traceID), func(ctx context.Context) {
				h.ServeHTTP(w, r.WithContext(ctx))
			})
			return
		}

		h.ServeHTTP(w, r)
	})
}
