package logkit

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Fields map[string]interface{}
type Field = zapcore.Field

type Logger = zap.Logger
type Core = zapcore.Core
type LKLevel = zapcore.Level

var (
	Any         = zap.Any
	Array       = zap.Array
	Binary      = zap.Binary
	Bool        = zap.Bool
	Bools       = zap.Bools
	ByteString  = zap.ByteString
	ByteStrings = zap.ByteStrings
	Complex128  = zap.Complex128
	Complex128s = zap.Complex128s
	Complex64   = zap.Complex64
	Complex64s  = zap.Complex64s
	Duration    = zap.Duration
	Durations   = zap.Durations
	Err         = zap.Error
	Errors      = zap.Errors
	Float32     = zap.Float32
	Float32s    = zap.Float32s
	Float64     = zap.Float64
	Float64s    = zap.Float64s
	Int         = zap.Int
	Int16       = zap.Int16
	Int16s      = zap.Int16s
	Int32       = zap.Int32
	Int32s      = zap.Int32s
	Int64       = zap.Int64
	Int64s      = zap.Int64s
	Int8        = zap.Int8
	Int8s       = zap.Int8s
	Ints        = zap.Ints
	NamedError  = zap.NamedError
	Namespace   = zap.Namespace
	Object      = zap.Object
	Reflect     = zap.Reflect
	Skip        = zap.Skip
	Stack       = zap.Stack
	String      = zap.String
	Stringer    = zap.Stringer
	Strings     = zap.Strings
	Time        = zap.Time
	Times       = zap.Times
	Uint        = zap.Uint
	Uint16      = zap.Uint16
	Uint16s     = zap.Uint16s
	Uint32      = zap.Uint32
	Uint32s     = zap.Uint32s
	Uint64      = zap.Uint64
	Uint64s     = zap.Uint64s
	Uint8       = zap.Uint8
	Uint8s      = zap.Uint8s
	Uintptr     = zap.Uintptr
	Uintptrs    = zap.Uintptrs
	Uints       = zap.Uints
)
