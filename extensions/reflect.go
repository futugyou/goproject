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

// Determines whether all fields are nil. This is only valid for simple objects.
func IsAllFieldsNil(v interface{}) bool {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := field.Type()

		switch field.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan:
			if !field.IsNil() {
				return false
			}
		case reflect.Struct:
			if !IsAllFieldsNil(field.Interface()) {
				return false
			}
		case reflect.Array:
			for j := 0; j < field.Len(); j++ {
				if !reflect.DeepEqual(field.Index(j).Interface(), reflect.Zero(fieldType.Elem()).Interface()) {
					return false
				}
			}
		default:
			if !reflect.DeepEqual(field.Interface(), reflect.Zero(fieldType).Interface()) {
				return false
			}
		}
	}
	return true
}
