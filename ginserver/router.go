package ginserver

import (
	"context"
	"fmt"
	"github.com/lee31802/comment_lib/constants"
	"github.com/lee31802/comment_lib/errors"
	"github.com/lee31802/comment_lib/logkit"
	"github.com/lee31802/comment_lib/trace"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/codegangsta/inject"

	"github.com/gin-gonic/gin"
)

// Handler can be any type of function.
type Handler interface{}

// Router is application router
type Router interface {
	Use(middlewares ...gin.HandlerFunc)
	Group(relativePath string, middlewares ...gin.HandlerFunc) Router
	Handle(method, path string, handler Handler, middlewares ...gin.HandlerFunc)
	GET(path string, handler Handler, middlewares ...gin.HandlerFunc)
	POST(path string, handler Handler, middlewares ...gin.HandlerFunc)
	DELETE(path string, handler Handler, middlewares ...gin.HandlerFunc)
	PATCH(path string, handler Handler, middlewares ...gin.HandlerFunc)
	PUT(path string, handler Handler, middlewares ...gin.HandlerFunc)
	OPTIONS(path string, handler Handler, middlewares ...gin.HandlerFunc)
	HEAD(path string, handler Handler, middlewares ...gin.HandlerFunc)
}

type router struct {
	path     string
	injector inject.Injector
	rg       *gin.RouterGroup
}

// Group creates a new router. You should add all the routes that have common middlwares or the same path prefix.
func (r *router) Group(relativePath string, middlewares ...gin.HandlerFunc) Router {
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}
	return &router{
		path:     r.path + relativePath,
		injector: r.injector,
		rg:       r.rg.Group(relativePath, middlewares...),
	}
}

// Handle registers a new request handle and middleware with the given path and method.
func (r *router) Handle(method, path string, handler Handler, middlewares ...gin.HandlerFunc) {
	addHandlerInfo(method, strings.Replace(r.path+path, "//", "/", -1), handler, middlewares)
	var chain []gin.HandlerFunc
	for _, middleware := range middlewares {
		chain = append(chain, middleware)
	}
	chain = append(chain, r.wraphandler(handler))
	r.rg.Handle(method, path, chain...)
}

// Use adds middleware to the group.
func (r *router) Use(middlewares ...gin.HandlerFunc) {
	r.rg.Use(middlewares...)
}

// GET is a shortcut for router.Handle("GET", path, handler, middlewares...).
func (r *router) GET(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("GET", path, handler, middlewares...)
}

// POST is a shortcut for router.Handle("POST", path, handler, middlewares...).
func (r *router) POST(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("POST", path, handler, middlewares...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handler, middlewares...).
func (r *router) DELETE(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("DELETE", path, handler, middlewares...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handler, middlewares...).
func (r *router) PATCH(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("PATCH", path, handler, middlewares...)
}

// PUT is a shortcut for router.Handle("PUT", path, handler, middlewares...).
func (r *router) PUT(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("PUT", path, handler, middlewares...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handler, middlewares...).
func (r *router) OPTIONS(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("OPTIONS", path, handler, middlewares...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handler, middlewares...).
func (r *router) HEAD(path string, handler Handler, middlewares ...gin.HandlerFunc) {
	r.Handle("HEAD", path, handler, middlewares...)
}

func (r *router) wraphandler(f Handler) gin.HandlerFunc {
	return convertHandler(f, r.injector)
}

func newReqInstance(t reflect.Type) interface{} {
	if t == nil {
		return nil
	}
	switch t.Kind() {
	case reflect.Ptr:
		return newReqInstance(t.Elem())
	case reflect.Interface:
		return nil
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0).Interface()
	case reflect.Array:
		v := reflect.New(t).Elem()
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			elementValue := reflect.ValueOf(newReqInstance(element.Type()))
			if elementValue.IsValid() && elementValue.Type().AssignableTo(element.Type()) {
				element.Set(elementValue)
			}
		}
		return v.Interface()
	case reflect.Map:
		m := reflect.MakeMap(t)
		return m.Interface()
	case reflect.String:
		return "default"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(0)
	case reflect.Float32, reflect.Float64:
		return float64(0)
	case reflect.Bool:
		return false
	default:
		return reflect.New(t).Interface()
	}
}

func newRespInstance(t reflect.Type) interface{} {
	if t == nil {
		return nil
	}
	fmt.Printf("%v\n", t.Kind())
	switch t.Kind() {
	case reflect.Ptr:
		return newRespInstance(t.Elem())
	case reflect.Interface:
		if t.NumMethod() == 0 {
			// 如果是没有方法的空接口，返回一个空结构体实例
			return struct{}{}
		}
		// 处理 Response 接口
		if t.Implements(reflect.TypeOf((*Response)(nil)).Elem()) {
			// 创建一个 jsonResponse 实例
			return &jsonResponse{}
		}
		// 这里可以根据具体接口要求创建更合适的实现，暂时返回 nil
		return nil
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0).Interface()
	case reflect.Array:
		v := reflect.New(t).Elem()
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			elementValue := reflect.ValueOf(newRespInstance(element.Type()))
			if elementValue.IsValid() && elementValue.Type().AssignableTo(element.Type()) {
				element.Set(elementValue)
			}
		}
		return v.Interface()
	case reflect.Map:
		m := reflect.MakeMap(t)
		return m.Interface()
	case reflect.String:
		return "default"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint64(0)
	case reflect.Float32, reflect.Float64:
		return float64(0)
	case reflect.Bool:
		return false
	default:
		return reflect.New(t).Interface()
	}
}

// convert to gin.HandlerFunc
func convertHandler(f Handler, parentInjector inject.Injector) gin.HandlerFunc {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic("handler should be a function")
	}
	// 返回值数量为 0：不做额外检查。
	// 返回值数量为 1：检查返回值类型是否为 string 或者实现了 Response 接口。如果不满足条件，触发 panic。
	// 返回值数量为 3：检查第一个返回值类型是否为 int(httpcode)，第二个返回值是否实现了 errors.Error 接口。如果不满足条件，触发 panic。
	// 其他返回值数量：触发 panic。
	switch t.NumOut() {
	case 0:
	case 1:
		outTyp := t.Out(0)
		if outTyp.Kind() != reflect.String && !outTyp.Implements(reflect.TypeOf((*Response)(nil)).Elem()) {
			panic("handler output parameter type should be `string` or `ginweb.Response`")
		}
	case 3:
		if codeTyp := t.Out(0); codeTyp.Kind() != reflect.Int {
			panic("handler first parameter type should be `int`")
		}
		if errTyp := t.Out(1); !errTyp.Implements(reflect.TypeOf((*errors.Error)(nil)).Elem()) {
			panic("handler second parameter type should be `gerrors.Error`")
		}
	default:
		panic("handler output parameter count should be 0, 1 or 3")
	}

	numIn := t.NumIn()
	requestFields := []int{}
	for i := 0; i < numIn; i++ {
		if t.In(i).Implements(reflect.TypeOf((*ServiceRequest)(nil)).Elem()) {
			requestFields = append(requestFields, i)
		}
	}
	if len(requestFields) > 1 {
		panic("handler should only have one request")
	}
	var handlerName string
	return func(c *gin.Context) {
		traceId := c.GetHeader(constants.CtxKeyTraceID)
		if traceId == "" {
			traceId = trace.NewRequestID()
		}
		rid := c.GetHeader(constants.CtxKeyRequestID)
		if rid == "" {
			rid = trace.NewRequestID()
			c.Set(constants.CtxKeyRequestID, rid)
		}
		traceCtx := trace.NewContextWithRequestID(c, rid)
		traceCtx = trace.NewContextWithTraceID(traceCtx, traceId)
		traceCtx = logkit.NewContextWith(traceCtx, trace.FieldTraceID(traceId), trace.FieldRequestID(rid))
		gwCtx := &constants.Context{
			Context: traceCtx,
			GinCtx:  c,
		}
		if handlerName == "" {
			handlerName = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		}
		c.Set(constants.CtxKeyHandlerName, handlerName)
		injector := inject.New()
		if parentInjector != nil {
			injector.SetParent(parentInjector)
		}
		for _, field := range requestFields {
			if req := newReqInstance(t.In(field)); req != nil {
				rp := newRequestParser(req)
				// 将req, path query等通过gin的context绑定起来
				if err := rp.parse(c); err != nil {
					logkit.FromContext(traceCtx).Error("Parse request failed", logkit.Err(err), logkit.Any("req", req))
					c.AbortWithStatusJSON(http.StatusBadRequest, &jsonResponseData{ErrCode: err.GetCode(), ErrMsg: err.GetMsg()})
					return
				}
				gr := req.(ServiceRequest)
				if err := gr.Parse(c); err != nil && err != errors.Success {
					logkit.FromContext(traceCtx).Error("Parse request failed", logkit.Err(err), logkit.Any("req", req))
					c.AbortWithStatusJSON(http.StatusBadRequest, &jsonResponseData{ErrCode: err.GetCode(), ErrMsg: err.GetMsg()})
					return
				}
				if err := gr.Validate(); err != nil && err != errors.Success {
					logkit.FromContext(traceCtx).Error("Validate request failed", logkit.Err(err), logkit.Any("req", req))
					c.AbortWithStatusJSON(http.StatusBadRequest, &jsonResponseData{ErrCode: err.GetCode(), ErrMsg: err.GetMsg()})
					return
				}
				injector.Map(req)
			}
		}
		injector.Map(gwCtx)
		injector.MapTo(gwCtx, (*context.Context)(nil))
		injector.Map(c)
		ret, err := injector.Invoke(f)
		if err != nil {
			panic(err)
		}

		switch len(ret) {
		case 1:
			i := ret[0].Interface()
			switch i.(type) {
			case Response:
				if i != nil {
					i.(Response).Render(c)
				}
			case string:
				c.String(http.StatusOK, i.(string))
			}
		case 3:
			gErr := ret[1].Interface().(errors.Error)
			resp := &jsonResponse{
				jsonResponseData: jsonResponseData{
					ErrCode: gErr.GetCode(),
					ErrMsg:  gErr.GetMsg(),
					Data:    ret[2].Interface(),
				},
				httpStatusCode: int(ret[0].Int()),
			}
			if c.Query("_show_request_id") != "" {
				resp.jsonResponseData.RequestID = &rid
			}
			resp.Render(c)
		}
	}
}

// ConvertHandler converts a ginweb handler to gin handler.
func ConvertHandler(f Handler) gin.HandlerFunc {
	return convertHandler(f, nil)
}
