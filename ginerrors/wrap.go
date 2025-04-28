package ginerrors

import (
	"fmt"
	"io"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Causer interface {
	Cause() error
	Params() []zapcore.Field
}

// Wrap returns an error annotating err the param.
// If err is nil, Wrap returns nil.
func Wrap(err error, params ...zapcore.Field) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause:  err,
		params: params,
	}
}

// Wrapf returns an error annotating err with the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause:  err,
		params: []zap.Field{zap.String("msg", fmt.Sprintf(format, args...))},
	}
}

type withMessage struct {
	cause  error
	params []zapcore.Field
}

func (w *withMessage) Error() string           { return w.cause.Error() }
func (w *withMessage) Cause() error            { return w.cause }
func (w *withMessage) Params() []zapcore.Field { return w.params }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// Unwrap returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the
// interface: Causer
//
// If the error does not implement Cause, it will return nil
// If the error is nil, nil will be returned without further
// investigation.
func Unwrap(err error) error {
	cause, ok := err.(Causer)
	if !ok {
		return nil
	}
	return cause.Cause()
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func As(err error, target interface{}) bool {
	if target == nil {
		panic("ginerrors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("ginerrors: target must be a non-nil pointer")
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && !e.Implements(errorType) {
		panic("ginerrors: *target must be interface or implement error")
	}
	targetType := typ.Elem()
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

// Params returns the input params cause the error, if possible.
// An error value has a cause&params if it implements the following
// interface: Causer
//
// If the error does not implement Cause, the original error will
// be returned. Without further investigation.
func Params(err error) []zapcore.Field {
	var fields []zapcore.Field
	for err != nil {
		cause, ok := err.(Causer)
		if !ok {
			break
		}
		fields = append(fields, cause.Params()...)
		err = cause.Cause()
	}
	return fields
}
