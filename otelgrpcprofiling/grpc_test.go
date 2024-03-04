package otelgrpcprofiling

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGrpcMiddlewareMatch(t *testing.T) {
	tests := []struct {
		name           string
		includePrefix  []string
		ignorePrefix   []string
		fullMethodName string
		expect         bool
	}{
		{
			name:           "Match any",
			includePrefix:  nil,
			ignorePrefix:   nil,
			fullMethodName: "/test.Prefix/Method",
			expect:         true,
		},
		{
			name:           "Included prefix",
			includePrefix:  []string{"/test.Prefix"},
			ignorePrefix:   nil,
			fullMethodName: "/test.Prefix/Method",
			expect:         true,
		},
		{
			name:           "Ignored prefix",
			includePrefix:  []string{"/test.Prefix"},
			ignorePrefix:   []string{"/test.Prefix/Ignore"},
			fullMethodName: "/test.Prefix/IgnoreMethod",
			expect:         false,
		},
		{
			name:           "No match",
			includePrefix:  []string{"/test.Prefix"},
			ignorePrefix:   nil,
			fullMethodName: "/other.Prefix/Method",
			expect:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{}
			for _, prefix := range tt.includePrefix {
				opts = append(opts, WithPrefix(prefix))
			}
			for _, prefix := range tt.ignorePrefix {
				opts = append(opts, WithIgnorePrefix(prefix))
			}
			middleware := NewMiddleware(opts...)

			require.Equal(t, tt.expect, middleware.match(tt.fullMethodName))
		})
	}
}
