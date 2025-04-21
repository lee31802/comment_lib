package server

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ctxbBinders = []contextBinder{
	pathBinding{},
	queryBinding{},
	headerBinding{},
}

type contextBinder interface {
	Name() string
	Bind(ctx *gin.Context, typeField *reflect.StructField, structField *reflect.Value) error
}

type queryBinding struct{}

func (pb queryBinding) Name() string { return "query" }
func (pb queryBinding) Bind(ctx *gin.Context, typeField *reflect.StructField, structField *reflect.Value) error {
	inputFieldName, ok := typeField.Tag.Lookup(pb.Name())
	if !ok {
		return nil
	}
	inputFieldNameList := strings.Split(inputFieldName, ",")
	inputFieldName = inputFieldNameList[0]
	var defaultValue string
	if len(inputFieldNameList) > 1 {
		defaultList := strings.SplitN(inputFieldNameList[1], "=", 2)
		if defaultList[0] == "default" {
			defaultValue = defaultList[1]
		}
	}
	if inputFieldName == "" {
		inputFieldName = typeField.Name
	}
	inputValue := ctx.Query(inputFieldName)
	if inputValue == "" && defaultValue != "" {
		inputValue = defaultValue
	}
	if inputValue == "" && structField.Kind() == reflect.Ptr {
		return nil
	}
	if _, isTime := structField.Interface().(time.Time); isTime {
		return setTimeField(inputValue, *typeField, *structField)
	}
	return setWithProperType(typeField.Type.Kind(), inputValue, *structField)
}

type pathBinding struct{}

func (pb pathBinding) Name() string { return "path" }
func (pb pathBinding) Bind(ctx *gin.Context, typeField *reflect.StructField, structField *reflect.Value) error {
	inputFieldName, ok := typeField.Tag.Lookup(pb.Name())
	if !ok {
		return nil
	}
	inputFieldNameList := strings.Split(inputFieldName, ",")
	inputFieldName = inputFieldNameList[0]
	var defaultValue string
	if len(inputFieldNameList) > 1 {
		defaultList := strings.SplitN(inputFieldNameList[1], "=", 2)
		if defaultList[0] == "default" {
			defaultValue = defaultList[1]
		}
	}
	if inputFieldName == "" {
		inputFieldName = typeField.Name
	}
	inputValue, exists := ctx.Params.Get(inputFieldName)
	if !exists && defaultValue != "" {
		inputValue = defaultValue
	}
	if inputValue == "" && structField.Kind() == reflect.Ptr {
		return nil
	}
	if _, isTime := structField.Interface().(time.Time); isTime {
		return setTimeField(inputValue, *typeField, *structField)
	}
	return setWithProperType(typeField.Type.Kind(), inputValue, *structField)
}

type headerBinding struct{}

func (hb headerBinding) Name() string { return "head" }
func (hb headerBinding) Bind(ctx *gin.Context, typeField *reflect.StructField, structField *reflect.Value) error {
	inputFieldName, ok := typeField.Tag.Lookup(hb.Name())
	if !ok {
		return nil
	}
	inputFieldNameList := strings.Split(inputFieldName, ",")
	inputFieldName = inputFieldNameList[0]
	var defaultValue string
	if len(inputFieldNameList) > 1 {
		defaultList := strings.SplitN(inputFieldNameList[1], "=", 2)
		if defaultList[0] == "default" {
			defaultValue = defaultList[1]
		}
	}
	if inputFieldName == "" {
		inputFieldName = typeField.Name
	}
	inputValue := ctx.GetHeader(inputFieldName)
	if inputValue == "" {
		inputValue = defaultValue
	}
	if inputValue == "" && structField.Kind() == reflect.Ptr {
		return nil
	}
	if _, isTime := structField.Interface().(time.Time); isTime {
		return setTimeField(inputValue, *typeField, *structField)
	}
	return setWithProperType(typeField.Type.Kind(), inputValue, *structField)
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	case reflect.Ptr:
		if !structField.Elem().IsValid() {
			structField.Set(reflect.New(structField.Type().Elem()))
		}
		structFieldElem := structField.Elem()
		return setWithProperType(structFieldElem.Kind(), val, structFieldElem)
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		return errors.New("Blank time format")
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
