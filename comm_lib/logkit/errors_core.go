package logkit

import (
	"github.com/lee31802/golog/comm_lib/errors"
	"go.uber.org/zap/zapcore"
)

func NewErrorsExtractCore(c Core) Core {
	return &errExtraCore{c}
}

type errExtraCore struct {
	zapcore.Core
}

func (c *errExtraCore) With(fields []Field) Core {
	fields = append(fields, extractFields(fields)...)
	return &errExtraCore{
		c.Core.With(fields),
	}
}

func extractFields(fields []zapcore.Field) []zapcore.Field {
	var retFields = make([]zapcore.Field, 0)
	for _, field := range fields {
		if field.Type != zapcore.ErrorType {
			continue
		}
		err := field.Interface
		for err != nil {
			cause, ok := err.(errors.Causer)
			if !ok {
				break
			}
			retFields = append(retFields, cause.Params()...)
			err = cause.Cause()
		}
	}
	return retFields
}
