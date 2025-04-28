package ginerrors

import (
	"encoding/json"
	"fmt"
)

type Error interface {
	error
	GetCode() int32
	GetMsg() string
}

type err GinError

func (e *err) Error() string {
	return fmt.Sprintf("error: code = %d desc = %s", e.GetCode(), e.GetMsg())
}

func (e *err) GetCode() int32 {
	return e.Code
}

func (e *err) GetMsg() string {
	return e.Msg
}

func New(code int32, msg string) Error {
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
