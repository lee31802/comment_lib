package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lee31802/comment_lib/errors"
	"reflect"
)

const (
	TraceIDHeader   = "Trace-ID"
	RequestIDHeader = "Request-ID"
)

// Request is ginweb's request structure that gets embedded in user defined request.
type Request struct{}

// Parse parses request from gin context.
func (r *Request) Parse(c *gin.Context) errors.Error {
	return errors.Success
}

// Validate checks the validation of the request.
func (r *Request) Validate() errors.Error {
	return errors.Success
}

type ServerRequest interface {
	Parse(*gin.Context) errors.Error
	Validate() errors.Error
}

type requestParser struct {
	req interface{}
	err error
}

func newRequestParser(req interface{}) *requestParser {
	return &requestParser{
		req: req,
	}
}

func (rp *requestParser) parse(c *gin.Context) errors.Error {
	rp.err = c.ShouldBind(rp.req)
	rp.bindContext(c, rp.req)
	if rp.err != nil {
		return errors.ErrorParseRequest
	}
	return nil
}

// 绑定上下文
func (rp *requestParser) bindContext(c *gin.Context, s interface{}) {
	typ := reflect.TypeOf(s)
	val := reflect.ValueOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			typeField := typ.Field(i)
			structField := val.Field(i)
			if !structField.CanSet() {
				continue
			}
			structFieldKind := structField.Kind()
			switch structFieldKind {
			case reflect.Ptr:
				if reflect.ValueOf(structField.Interface()).Elem().IsValid() {
					//type Address struct {
					//    Street string
					//    City   string
					//}
					//
					//type User struct {
					//    Name    string
					//    Address *Address
					//}
					// 针对address情况就需要进行递归调用
					rp.bindContext(c, structField.Interface())
					continue
				}
				// v := reflect.ValueOf(newReqInstance(structField.Type()))
				// if v.Elem().Kind() == reflect.Struct {
				// 	structField.Set(v)
				// 	rp.bindContext(c, structField.Interface())
				// 	continue
				// }
			case reflect.Struct:
				//type Address struct {
				//    Street string
				//    City   string
				//}
				//
				//type User struct {
				//    Name    string
				//    Address *Address
				//}
				// 针对address情况就需要进行递归调用
				rp.bindContext(c, structField.Addr().Interface())
				continue
			}
			for _, binder := range ctxbBinders {
				err := binder.Bind(c, &typeField, &structField)
				if err != nil {
					rp.err = err
					return
				}
			}
		}
	}
}
