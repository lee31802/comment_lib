package trace

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type key int

// key list
const (
	KeyRequestID  key = iota
	KeyStressFlag key = iota
	KeyTraceID    key = iota
)

// NewContextWithRequestID injects a request id string to context.
func NewContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, KeyRequestID, requestID)
}

func NewContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, KeyTraceID, traceID)
}

func NewContextWithMarkStressFlag(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeyStressFlag, true)
}

func IsStressContext(ctx context.Context) bool {
	stressFlag, ok := ctx.Value(KeyStressFlag).(bool)
	return ok && stressFlag
}

// RequestIDFromContext extracts a request id string from context.
func RequestIDFromContext(ctx context.Context) string {
	rid, ok := ctx.Value(KeyRequestID).(string)
	if !ok {
		return ""
	}
	return rid
}

// NewRequestID returns a unique request id string.
func NewRequestID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
