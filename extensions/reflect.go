package extensions

import (
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
		if val == nil {
			return nil
		}
		return val
	}
	if v.Kind() == reflect.String {
		val := v.String()
		if len(val) == 0 {
			return nil
		}
		return &val
	}

	return nil
}

func GetStringFieldStruct(obj interface{}, fields ...string) string {
	v := GetStringFieldPointer(obj, fields...)
	if v == nil {
		return ""
	}
	return *v
}
