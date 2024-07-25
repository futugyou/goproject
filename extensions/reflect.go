package extensions

import (
	"fmt"
	"reflect"
)

func GetStringFieldPointer(obj interface{}, fields ...string) *string {
	v := reflect.ValueOf(obj)
	for _, field := range fields {
		if v.Kind() == reflect.Ptr && !v.IsNil() {
			v = v.Elem()
		} else {
			return nil
		}
		v = v.FieldByName(field)
		if !v.IsValid() {
			return nil
		}
	}
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		val := v.Interface().(*string)
		return val
	}
	if v.Kind() == reflect.String {
		val := v.String()
		return &val
	}
	fmt.Println(v.Kind() == reflect.String)
	return nil
}

func GetStringFieldStruct(obj interface{}, fields ...string) string {
	v := GetStringFieldPointer(obj, fields...)
	if v == nil {
		return ""
	}
	return *v
}
