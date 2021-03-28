package correlation

import (
	"context"
)

type correlationContextKeyType string

const correlationContextKey correlationContextKeyType = "correlation"

func ContextWithCorrelation(parent context.Context, correlation string) context.Context {
	return context.WithValue(parent, correlationContextKey, correlation)
}

func CorrelationFromContext(ctx context.Context) string {
	correlation, ok := ctx.Value(correlationContextKey).(string)
	if !ok {
		return "UNDEFINED"
	}
	return correlation
}
