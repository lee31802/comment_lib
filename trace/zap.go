package trace

import "go.uber.org/zap"

// FieldRequestID returns a zap.Field of request id.
func FieldRequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

func FieldTraceID(traceID string) zap.Field {
	return zap.String("trace_id", traceID)
}
