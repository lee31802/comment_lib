package constants

import (
	"context"
	"github.com/gin-gonic/gin"
)

var (
	CtxKeyRequestID   = "request_id"
	CtxKeyTraceID     = "trace_id"
	CtxKeyHandlerName = "handler_name"
)

type Context struct {
	context.Context
	GinCtx *gin.Context
}
