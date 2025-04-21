package errors

import (
	"encoding/json"
	"fmt"
)

var (
	//common error : [0,999]
	Success           = &BaseError{0, "success"}
	ErrorUnkown       = &BaseError{1, "unknown"}
	ErrorParseRequest = &BaseError{2, "parse request error"}

	ErrorParamsInvalid    = &BaseError{100, "params invalid"}
	ErrorDuplicateRequest = &BaseError{101, "duplicate request"}
	ErrorUserSigCal       = &BaseError{102, "usersig cal error"}
	ErrorJSONMarshal      = &BaseError{103, "json marshal error"}
	ErrorJSONUnMarshal    = &BaseError{104, "json unmarshal error"}

	ErrorDBInit       = &BaseError{200, "database init error"}
	ErrorRedisInit    = &BaseError{201, "redis init error"}
	ErrorCacheInit    = &BaseError{202, "cache init error"}
	ErrorDBOperate    = &BaseError{203, "database operate error"}
	ErrorRedisOperate = &BaseError{204, "redis operate error"}
	ErrorCacheOperate = &BaseError{205, "cache operate error"}
	ErrorDBTxBegin    = &BaseError{206, "database tx begin error"}
	ErrorDBTxCommit   = &BaseError{207, "database tx commit error"}
	ErrorDataNotFound = &BaseError{208, "data not found"}

	ErrorRPCCall     = &BaseError{300, "rpc call error"}
	ErrorCallTimeout = &BaseError{301, "call timeout error"}
	ErrorBrokenPipe  = &BaseError{302, "broken pipe"}

	//others
	ErrorUnknown           = &BaseError{500, "server internal error"}
	ErrorInvalidValidation = &BaseError{501, "Invalid Validation Error"}
	ErrorValidation        = &BaseError{502, "Validation Error"}
)

type Error interface {
	error
	GetCode() int32
	GetMsg() string
}

type err PbRPCError

func (e err) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s", e.GetCode(), e.GetMsg())
}

func (e err) GetCode() int32 {
	return e.Code
}

func (e err) GetMsg() string {
	return e.Msg
}

func New(code int32, msg string) error {
	return &err{
		Code: code,
		Msg:  msg,
	}
}

func Errorf(code int32, format string, params ...interface{}) error {
	return New(code, fmt.Sprintf(format, params...))
}

func Parse(origin error) (error, bool) {
	if e, ok := origin.(Error); ok {
		return New(e.GetCode(), e.GetMsg()), true
	}
	return New(0, origin.Error()), false
}

func Json(origin string) (error, bool) {
	error := &err{}
	if e := json.Unmarshal([]byte(origin), error); e == nil {
		return error, true
	}
	return New(0, origin), false
}
