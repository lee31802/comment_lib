package server

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	globalHandlerInfos handlerInfoList
	// 主要用来全链路跟踪，记录每个节点的执行时间等
	pathHandlerMap map[string]string
)

type handlerInfoList []*handlerInfo

func init() {
	pathHandlerMap = make(map[string]string)
}

type handlerInfo struct {
	HandlerName string
	URL         string
	Method      string
	Request     *requestInfo
	Response    *responseInfo
}

type requestInfo struct {
	Name       string
	PkgPath    string
	CurlString string
	FieldInfos []*requestFieldInfo
}

type responseInfo struct {
	Desc       string
	FieldInfos []*responseFieldInfo
}

type requestFieldInfo struct {
	Name     string
	Typ      string
	Tag      string
	Required bool
}

type responseFieldInfo struct {
	Name string
	Typ  string
	Tag  string
}

// simplify handlerName, only keep module name and method name
// for example, the raw handlerName is examples/blog/article.(*ArticleModule).GetArticles-fm
// then,simplify it, we will get article.GetArticles
func getHandlerSimpleName(handlerName string) string {
	n := strings.LastIndexByte(handlerName, '/')
	handlerName = handlerName[n+1:]
	simpleName := strings.TrimSuffix(handlerName, "-fm")
	parts := strings.Split(simpleName, ".")
	if len(parts) == 3 {
		simpleName = fmt.Sprintf("%s.%s", parts[0], parts[2])
	}
	return simpleName
}

// addHandlerInfo 函数通过反射机制，收集和整理了处理函数的名称、请求方法、请求路径、输入参数和输出响应的相关信息，并将这些信息存储到 handlerInfo 结构体中，同时更新了一些全局的映射和列表。
func addHandlerInfo(method, path string, handler Handler, middlewares []gin.HandlerFunc) *handlerInfo {
	ht := reflect.TypeOf(handler)
	withResponse := false
	switch ht.NumOut() {
	case 1:
	case 3:
		withResponse = true
	}
	// main.(*testModule).testFunc-fm
	handlerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	s := strings.Split(handlerName, ".")
	realHandlerName := strings.TrimRight(s[len(s)-1], "-fm")
	info := &handlerInfo{
		HandlerName: realHandlerName,
		URL:         path,
		Method:      method,
	}
	pathHandlerMap[info.Method+":"+info.URL] = getHandlerSimpleName(handlerName)
	t := reflect.TypeOf(handler)
	for i := 0; i < t.NumIn(); i++ {
		if t.In(i).Implements(reflect.TypeOf((*ServerRequest)(nil)).Elem()) {
			if req := newReqInstance(t.In(i)); req != nil {
				if withResponse {
					response := &responseInfo{}
					responseValue := reflect.New(ht.Out(2))
					if responseValue.Kind() == reflect.Ptr && reflect.TypeOf(responseValue.Elem().Interface()) != nil {
						responseValue = responseValue.Elem()
					}
					responseTyp := reflect.TypeOf(responseValue.Interface())
					if responseTyp.Kind() == reflect.Ptr {
						responseTyp = responseTyp.Elem()
					}
					if responseTyp.Kind() == reflect.Struct {
						var infos []*responseFieldInfo
						for j := 0; j < responseTyp.NumField(); j++ {
							typeField := responseTyp.Field(j)
							infos = append(infos, &responseFieldInfo{
								Name: typeField.Name,
								Typ:  typeField.Type.String(),
								Tag:  fmt.Sprintf("%v", typeField.Tag),
							})
						}
						response.FieldInfos = infos
					}
					info.Response = response
				}
				reqTyp := reflect.TypeOf(req).Elem()
				request := parseRequestTypeFields(reqTyp, method, path)
				info.Request = request
			}
		}
	}
	globalHandlerInfos = append(globalHandlerInfos, info)
	return info
}

func parseRequestTypeFields(t reflect.Type, method string, p string) *requestInfo {
	host := "127.0.0.1"
	if port := os.Getenv("PORT"); port != "" {
		host += ":" + port
	}
	path := path.Join(g.opts.RootPath, p)
	url := url.URL{
		Scheme: "http",
		Host:   host,
	}
	var fieldInfos []*requestFieldInfo
	jsons := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		if typeField.Type.String() == "ginweb.Request" {
			continue
		}
		info := &requestFieldInfo{
			Name: typeField.Name,
			Typ:  typeField.Type.String(),
			Tag:  fmt.Sprintf("%v", typeField.Tag),
		}
		if val, ok := typeField.Tag.Lookup("path"); ok {
			path = strings.Replace(path, ":"+val, "1", -1)
			info.Required = true
		} else if val, ok := typeField.Tag.Lookup("json"); ok {
			jsons[val] = newReqInstance(typeField.Type)
			if val, ok := typeField.Tag.Lookup("binding"); ok && val == "required" {
				info.Required = true
			}
		} else if val, ok := typeField.Tag.Lookup("query"); ok {
			q := url.Query()
			q.Set(val, "1")
			url.RawQuery = q.Encode()
		}
		fieldInfos = append(fieldInfos, info)
	}
	url.Path = path
	curlString := fmt.Sprintf("curl -X%v '%v'", method, url.String())
	if len(jsons) > 0 {
		if buf, err := json.Marshal(&jsons); err == nil {
			curlString += fmt.Sprintf(" -H 'Content-Type: application/json' -d '%v'", string(buf))
		}
	}
	return &requestInfo{
		Name:       t.Name(),
		PkgPath:    t.PkgPath(),
		CurlString: curlString,
		FieldInfos: fieldInfos,
	}
}
