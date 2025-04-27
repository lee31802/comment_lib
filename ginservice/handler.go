package ginservice

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
	// api
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
	FieldInfos []*responseFieldInfo
}

type requestFieldInfo struct {
	Name        string
	Typ         string
	Tag         string
	Required    bool
	Description string
}

type responseFieldInfo struct {
	Name        string
	Typ         string
	Tag         string
	Description string
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
	fmt.Printf(":%v", t.NumIn())

	// 处理 request 信息
	for i := 0; i < t.NumIn(); i++ {
		if t.In(i).Implements(reflect.TypeOf((*ServiceRequest)(nil)).Elem()) {
			// 根据传入的反射类型 t 创建一个对应的实例
			if req := newReqInstance(t.In(i)); req != nil {
				reqTyp := reflect.TypeOf(req).Elem()
				request := parseRequestTypeFields(reqTyp, method, path)
				info.Request = request
			}
		}
	}
	fmt.Printf("%v", t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		if t.Out(i).Implements(reflect.TypeOf((*Response)(nil)).Elem()) {
			// 根据传入的反射类型 t 创建一个对应的实例
			if resp := newRespInstance(t.Out(i)); resp != nil {
				respTyp := reflect.TypeOf(resp).Elem()
				response := parseResponseTypeFields(respTyp)
				info.Response = response
			}
		}
	}

	globalHandlerInfos = append(globalHandlerInfos, info)
	return info
}

func parseResponseTypeFields(t reflect.Type) *responseInfo {
	var fieldInfos []*responseFieldInfo
	jsons := make(map[string]interface{})
	var buildJSON func(t reflect.Type, parentKey string, jsons map[string]interface{})
	buildJSON = func(t reflect.Type, parentKey string, jsons map[string]interface{}) {
		for i := 0; i < t.NumField(); i++ {
			typeField := t.Field(i)
			key := typeField.Name
			if parentKey != "" {
				key = parentKey + "." + key
			}
			info := &responseFieldInfo{
				Name: key,
				Typ:  typeField.Type.String(),
				Tag:  fmt.Sprintf("%v", typeField.Tag),
			}
			if val, ok := typeField.Tag.Lookup("json"); ok {
				if typeField.Type.Kind() == reflect.Struct {
					subJsons := make(map[string]interface{})
					buildJSON(typeField.Type, key, subJsons)
					jsons[val] = subJsons
				} else {
					jsons[val] = newReqInstance(typeField.Type)
				}
			}
			// 提取 desc 信息
			if desc, ok := typeField.Tag.Lookup("desc"); ok {
				info.Description = desc
			}
			fieldInfos = append(fieldInfos, info)
		}
	}
	buildJSON(t, "", jsons)
	return &responseInfo{
		FieldInfos: fieldInfos,
	}
}

func parseRequestTypeFields(t reflect.Type, method string, p string) *requestInfo {
	serviceHost := "127.0.0.1"
	if servicePort := os.Getenv("PORT"); servicePort != "" {
		serviceHost += ":" + servicePort
	}
	servicePath := path.Join(s.opts.RootPath, p)
	serviceUrl := url.URL{
		Scheme: "http",
		Host:   serviceHost,
	}
	var fieldInfos []*requestFieldInfo
	jsons := make(map[string]interface{})

	var buildJSON func(t reflect.Type, parentKey string, jsons map[string]interface{})
	buildJSON = func(t reflect.Type, parentKey string, jsons map[string]interface{}) {
		for i := 0; i < t.NumField(); i++ {
			typeField := t.Field(i)
			if typeField.Type.String() == "ginservice.Request" {
				continue
			}
			key := typeField.Name
			if parentKey != "" {
				key = parentKey + "." + key
			}
			info := &requestFieldInfo{
				Name: key,
				Typ:  typeField.Type.String(),
				Tag:  fmt.Sprintf("%v", typeField.Tag),
			}
			if val, ok := typeField.Tag.Lookup("path"); ok {
				servicePath = strings.Replace(servicePath, ":"+val, "1", -1)
				info.Required = true
			} else if val, ok := typeField.Tag.Lookup("json"); ok {
				if typeField.Type.Kind() == reflect.Struct {
					subJsons := make(map[string]interface{})
					buildJSON(typeField.Type, key, subJsons)
					jsons[val] = subJsons
				} else {
					jsons[val] = newReqInstance(typeField.Type)
				}
				if val, ok = typeField.Tag.Lookup("binding"); ok && val == "required" {
					info.Required = true
				}
			} else if val, ok := typeField.Tag.Lookup("query"); ok {
				q := serviceUrl.Query()
				q.Set(val, "1")
				serviceUrl.RawQuery = q.Encode()
			}
			// 提取 desc 信息
			if desc, ok := typeField.Tag.Lookup("desc"); ok {
				info.Description = desc
			}
			fieldInfos = append(fieldInfos, info)
		}
	}

	buildJSON(t, "", jsons)

	serviceUrl.Path = servicePath
	curlString := fmt.Sprintf("curl -X%v '%v'", method, serviceUrl.String())
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

//func parseRequestTypeFields(t reflect.Type, method string, p string) *requestInfo {
//	serviceHost := "127.0.0.1"
//	if servicePort := os.Getenv("PORT"); servicePort != "" {
//		serviceHost += ":" + servicePort
//	}
//	servicePath := path.Join(s.opts.RootPath, p)
//	serviceUrl := url.URL{
//		Scheme: "http",
//		Host:   serviceHost,
//	}
//	var fieldInfos []*requestFieldInfo
//	jsons := make(map[string]interface{})
//	for i := 0; i < t.NumField(); i++ {
//		typeField := t.Field(i)
//		if typeField.Type.String() == "ginservice.Request" {
//			continue
//		}
//		info := &requestFieldInfo{
//			Name: typeField.Name,
//			Typ:  typeField.Type.String(),
//			Tag:  fmt.Sprintf("%v", typeField.Tag),
//		}
//		if val, ok := typeField.Tag.Lookup("path"); ok {
//			// 构建一个可执行的curl命令
//			servicePath = strings.Replace(servicePath, ":"+val, "1", -1)
//			info.Required = true
//		} else if val, ok = typeField.Tag.Lookup("json"); ok {
//			jsons[val] = newReqInstance(typeField.Type)
//			if val, ok = typeField.Tag.Lookup("binding"); ok && val == "required" {
//				info.Required = true
//			}
//		} else if val, ok = typeField.Tag.Lookup("query"); ok {
//			q := serviceUrl.Query()
//			q.Set(val, "1")
//			serviceUrl.RawQuery = q.Encode()
//		}
//		fieldInfos = append(fieldInfos, info)
//	}
//	serviceUrl.Path = servicePath
//	curlString := fmt.Sprintf("curl -X%v '%v'", method, serviceUrl.String())
//	if len(jsons) > 0 {
//		if buf, err := json.Marshal(&jsons); err == nil {
//			curlString += fmt.Sprintf(" -H 'Content-Type: application/json' -d '%v'", string(buf))
//		}
//	}
//	return &requestInfo{
//		Name:       t.Name(),
//		PkgPath:    t.PkgPath(),
//		CurlString: curlString,
//		FieldInfos: fieldInfos,
//	}
//}
