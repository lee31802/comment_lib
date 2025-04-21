package encoder

import (
	"encoding/base64"
	"encoding/json"
	"github.com/lee31802/comment_lib/constants"
	"math"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// For JSON-escaping; see jsonEncoder.safeAddString below.
const (
	_hex = "0123456789abcdef"
)

type pool struct {
	_pool sync.Pool
}

func (pool *pool) get() *logEncoder {
	return pool._pool.Get().(*logEncoder)
}

func (pool *pool) put(enc *logEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.requestID = ""
	enc.traceID = ""
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	pool._pool.Put(enc)
}

var (
	bufferpool  = buffer.NewPool()
	encoderPool = &pool{
		_pool: sync.Pool{New: func() interface{} {
			return &logEncoder{}
		}},
	}
)

type logEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int
	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc *json.Encoder

	requestID string
	traceID   string
}

func NewEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &logEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
		spaced:        false,
	}
}

func (enc *logEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *logEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *logEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *logEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *logEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *logEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *logEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *logEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *logEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *logEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferpool.Get()
		enc.reflectEnc = json.NewEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

func (enc *logEncoder) AddReflected(key string, obj interface{}) error {
	enc.resetReflectBuf()
	err := enc.reflectEnc.Encode(obj)
	if err != nil {
		return err
	}
	enc.reflectBuf.TrimNewline()
	enc.addKey(key)
	_, err = enc.buf.Write(enc.reflectBuf.Bytes())
	return err
}

func (enc *logEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *logEncoder) AddString(key, val string) {
	switch key {
	case constants.KeyRequestID:
		enc.requestID = val
	case constants.KeyTraceID:
		enc.traceID = val
	default:
		enc.addKey(key)
		enc.AppendString(val)
	}
}

func (enc *logEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *logEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *logEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *logEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	return err
}

func (enc *logEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *logEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

func (enc *logEncoder) AppendComplex128(val complex128) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *logEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *logEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *logEncoder) AppendReflected(val interface{}) error {
	enc.resetReflectBuf()
	err := enc.reflectEnc.Encode(val)
	if err != nil {
		return err
	}
	enc.reflectBuf.TrimNewline()
	enc.addElementSeparator()
	_, err = enc.buf.Write(enc.reflectBuf.Bytes())
	return err
}

func (enc *logEncoder) AppendString(val string) {
	enc.addElementSeparator()
	enc.safeAddString(val)
}

func (enc *logEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *logEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *logEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *logEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *logEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *logEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *logEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *logEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *logEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *logEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *logEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *logEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *logEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *logEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *logEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *logEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *logEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *logEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *logEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *logEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *logEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *logEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *logEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *logEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *logEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *logEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *logEncoder) clone() *logEncoder {
	clone := encoderPool.get()
	clone.requestID = enc.requestID
	clone.traceID = enc.traceID
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferpool.Get()
	return clone
}

func (enc *logEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	addTimestamp(final, &ent)
	addLogLevel(final, &ent)
	addRequestID(final)
	addTraceID(final)
	addLoggerName(final, &ent)

	addCaller(final, &ent)
	addMsg(final, &ent)
	addBuffer(enc, final)
	addFields(final, fields)
	final.closeOpenNamespaces()
	addStacktrace(final, &ent)
	addLineEnding(final)
	ret := final.buf
	encoderPool.put(final)
	return ret, nil
}

func (enc *logEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *logEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
}

func (enc *logEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.safeAddString(key)
	enc.buf.AppendByte('=')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}

func (enc *logEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '|', '=', '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.AppendByte(',')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *logEncoder) addLogSegmentSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	if enc.buf.Bytes()[last] == '|' {
		return
	}

	enc.buf.AppendByte('|')
}

func (enc *logEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's logEncoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *logEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *logEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *logEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *logEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addTimestamp(enc *logEncoder, ent *zapcore.Entry) {
	if enc == nil || ent == nil {
		return
	}
	if enc.TimeKey != "" {
		enc.buf.AppendString(ent.Time.Format("2006-01-02 15:04:05.000"))
	}
}
func addLogLevel(enc *logEncoder, ent *zapcore.Entry) {
	if enc == nil || ent == nil {
		return
	}
	if enc.LevelKey != "" {
		enc.addLogSegmentSeparator()
		enc.buf.AppendString(strings.ToUpper(ent.Level.String()))
	}
}

func addLoggerName(enc *logEncoder, ent *zapcore.Entry) {
	if ent.LoggerName != "" && enc.NameKey != "" {
		enc.addKey(enc.NameKey)
		cur := enc.buf.Len()
		nameEncoder := enc.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, enc)
		if cur == enc.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output JSON valid.
			enc.AppendString(ent.LoggerName)
		}
	}
}

func addMsg(enc *logEncoder, ent *zapcore.Entry) {
	if enc.MessageKey != "" {
		enc.addLogSegmentSeparator()
		enc.AppendString(ent.Message)
	}
}

func addCaller(enc *logEncoder, ent *zapcore.Entry) {
	if ent.Caller.Defined && enc.CallerKey != "" {
		enc.addLogSegmentSeparator()
		cur := enc.buf.Len()
		enc.EncodeCaller(ent.Caller, enc)
		if cur == enc.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep output JSON valid.
			enc.AppendString(ent.Caller.String())
		}
		return
	}
	enc.buf.AppendString("|-")

}

func addBuffer(enc, final *logEncoder) {
	if enc.buf.Len() > 0 {
		final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes())
	}
}

// func add
func addRequestID(enc *logEncoder) {
	var requestID = "-"
	if enc.requestID != "" {
		requestID = enc.requestID
	}
	enc.addLogSegmentSeparator()
	enc.buf.AppendString(requestID)
}

// func add
func addTraceID(enc *logEncoder) {
	var traceID = "-"
	if enc.traceID != "" {
		traceID = enc.traceID
	}
	enc.addLogSegmentSeparator()
	enc.buf.AppendString(traceID)
}

// func add
func addFields(enc *logEncoder, fields []zapcore.Field) {
	if len(fields) == 0 {
		return
	}
	enc.addLogSegmentSeparator()
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
func addStacktrace(enc *logEncoder, ent *zapcore.Entry) {
	if ent.Stack != "" && enc.StacktraceKey != "" {
		enc.AddString(enc.StacktraceKey, ent.Stack)
	}
}
func addLineEnding(enc *logEncoder) {
	if enc.LineEnding != "" {
		enc.buf.AppendString(enc.LineEnding)
	} else {
		enc.buf.AppendString(zapcore.DefaultLineEnding)
	}
}
