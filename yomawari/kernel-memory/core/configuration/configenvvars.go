package configuration

import (
	"fmt"
	"reflect"
	"strings"
)

func GenerateEnvVarsFromObject(source interface{}, parents ...string) map[string]string {
	if source == nil {
		return make(map[string]string)
	}
	prefix := getPrefix(parents)
	return generateEnvVars(source, prefix)
}

func GenerateEnvVarsFromObjectNoDefaults(source interface{}, parents ...string) map[string]string {
	if source == nil {
		return make(map[string]string)
	}

	variables := GenerateEnvVarsFromObject(source, parents...)
	prefix := getPrefix(parents)

	defaults := generateEnvVars(createInstanceOfSameType(source), prefix)
	for k, v := range defaults {
		if val, exists := variables[k]; exists && val == v {
			delete(variables, k)
		}
	}

	return variables
}

func getPrefix(parents []string) string {
	if len(parents) > 0 {
		return strings.Join(parents, "__") + "__"
	}
	return ""
}

func generateEnvVars(source interface{}, prefix string) map[string]string {
	result := make(map[string]string)
	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		if !fieldValue.CanInterface() {
			continue
		}

		fullKey := fmt.Sprintf("%s%s", prefix, field.Name)
		switch fieldValue.Kind() {
		case reflect.Map:
			for _, key := range fieldValue.MapKeys() {
				mapValue := fieldValue.MapIndex(key)
				if mapValue.IsValid() && mapValue.CanInterface() {
					dictKey := fmt.Sprintf("%s__%v", fullKey, key)
					result[dictKey] = fmt.Sprintf("%v", mapValue.Interface())
				}
			}
		case reflect.Slice, reflect.Array:
			for index := 0; index < fieldValue.Len(); index++ {
				arrayKey := fmt.Sprintf("%s__%d", fullKey, index)
				result[arrayKey] = fmt.Sprintf("%v", fieldValue.Index(index).Interface())
			}
		default:
			result[fullKey] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	}
	return result
}

func createInstanceOfSameType(source interface{}) interface{} {
	typ := reflect.TypeOf(source)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return reflect.New(typ).Elem().Interface()
}
