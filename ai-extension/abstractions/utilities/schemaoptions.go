package utilities

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

const (
	SchemaPropertyName       = "$schema"
	TitlePropertyName        = "title"
	DescriptionPropertyName  = "description"
	TypePropertyName         = "type"
	PropertiesPropertyName   = "properties"
	RequiredPropertyName     = "required"
	AdditionalPropertiesName = "additionalProperties"
	SchemaKeywordUri         = "http://json-schema.org/draft-07/schema#"
)

// AIJsonSchemaCreateOptions defines configuration options when generating JSON Schema
type AIJsonSchemaCreateOptions struct {
	IncludeSchemaKeyword         bool
	IncludeTypeInEnumSchemas     bool
	DisallowAdditionalProperties bool
	RequireAllProperties         bool
}

// DefaultAIJsonSchemaCreateOptions default options
var DefaultAIJsonSchemaCreateOptions = &AIJsonSchemaCreateOptions{
	IncludeSchemaKeyword:         true,
	IncludeTypeInEnumSchemas:     false,
	DisallowAdditionalProperties: true,
	RequireAllProperties:         true,
}

// CreateFunctionJsonSchema generates a JSON Schema based on a function type, title, description, list of parameter names, and inferred options
func CreateFunctionJsonSchema(fnType reflect.Type, title string, description string, paramNames []string, inferenceOptions *AIJsonSchemaCreateOptions) (map[string]interface{}, error) {
	if fnType == nil {
		return nil, fmt.Errorf("type cannot be nil")
	}

	if inferenceOptions == nil {
		inferenceOptions = DefaultAIJsonSchemaCreateOptions
	}

	properties := make(map[string]interface{})
	required := make([]string, 0)

	numIn := fnType.NumIn()
	// Iterate over function arguments (skipping arguments that implement context.Context)
	for i := 0; i < numIn; i++ {
		inType := fnType.In(i)
		// If the parameter implements context.Context, skip
		if inType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			continue
		}
		// Parameter name: user-provided or automatically generated
		var paramName string
		if i < len(paramNames) {
			paramName = paramNames[i]
		} else {
			paramName = fmt.Sprintf("param%d", i)
		}

		paramSchema := createJsonSchemaCore(inType, paramName, "", nil, inferenceOptions)
		properties[paramName] = paramSchema
		required = append(required, paramName)
	}

	schema := make(map[string]interface{})
	if inferenceOptions.IncludeSchemaKeyword {
		schema[SchemaPropertyName] = SchemaKeywordUri
	}
	if title != "" {
		schema[TitlePropertyName] = title
	}
	if description != "" {
		schema[DescriptionPropertyName] = description
	}
	schema[TypePropertyName] = "object"
	schema[PropertiesPropertyName] = properties
	if len(required) > 0 {
		schema[RequiredPropertyName] = required
	}
	if inferenceOptions.DisallowAdditionalProperties {
		schema[AdditionalPropertiesName] = false
	}

	return schema, nil
}

// CreateJsonSchemaCore generates JSON Schema based on parameter type, name, description, etc. (simplified implementation)
func createJsonSchemaCore(
	t reflect.Type,
	paramName string,
	description string,
	defaultVal interface{},
	options *AIJsonSchemaCreateOptions,
) map[string]interface{} {

	schema := make(map[string]interface{})

	// base type
	switch t.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Struct:
		handleStructType(t, schema, options)
	case reflect.Ptr:
		handlePointerType(t, schema, paramName, description, defaultVal, options)
	case reflect.Slice, reflect.Array:
		handleArrayType(t, schema, paramName, options)
	case reflect.Map:
		handleMapType(schema, options)
	}

	// add description
	if description != "" {
		schema[DescriptionPropertyName] = description
	}

	// add default value
	if defaultVal != nil {
		schema["default"] = defaultVal
	}

	// handle enum values
	if enumValues := getEnumValues(t); len(enumValues) > 0 {
		schema["enum"] = enumValues
		if options.IncludeTypeInEnumSchemas {
			if _, exists := schema["type"]; !exists {
				schema["type"] = detectEnumType(enumValues)
			}
		}
	}

	return schema
}

// for go unused warning
var _ = parseJSONTag

func parseJSONTag(tag string) (name string, opts string) {
	parts := strings.Split(tag, ",")
	if len(parts) > 0 {
		name = parts[0]
	}
	if len(parts) > 1 {
		opts = strings.Join(parts[1:], ",")
	}
	return
}

// for go unused warning
var _ = parseDefaultValue

func parseDefaultValue(valStr string, t reflect.Type) interface{} {
	if valStr == "" {
		return nil
	}

	// TODO
	switch t.Kind() {
	case reflect.String:
		return valStr
	case reflect.Int:
		var i int
		fmt.Sscanf(valStr, "%d", &i)
		return i
	case reflect.Bool:
		return strings.ToLower(valStr) == "true"
	}
	return nil
}

func handleStructType(t reflect.Type, schema map[string]interface{}, options *AIJsonSchemaCreateOptions) {
	nestedSchema, _ := CreateFunctionJsonSchema(t, "", "", nil, options)
	schema["type"] = "object"
	schema["properties"] = nestedSchema["properties"]
	if req, ok := nestedSchema["required"]; ok {
		schema["required"] = req
	}
	if options.DisallowAdditionalProperties {
		schema["additionalProperties"] = false
	}
}

func handlePointerType(t reflect.Type, schema map[string]interface{}, paramName, description string, defaultVal interface{}, options *AIJsonSchemaCreateOptions) {
	elemType := t.Elem()
	elemSchema := createJsonSchemaCore(elemType, paramName, description, defaultVal, options)
	schema["anyOf"] = []map[string]interface{}{
		elemSchema,
		{"type": "null"},
	}
}

func handleArrayType(t reflect.Type, schema map[string]interface{}, paramName string, options *AIJsonSchemaCreateOptions) {
	elemType := t.Elem()
	elemSchema := createJsonSchemaCore(elemType, paramName, "", nil, options)
	schema["type"] = "array"
	schema["items"] = elemSchema
}

func handleMapType(schema map[string]interface{}, _ *AIJsonSchemaCreateOptions) {
	schema["type"] = "object"
	schema["additionalProperties"] = true
}

func getEnumValues(_ reflect.Type) []interface{} {
	// TODO
	return nil
}

func detectEnumType(values []interface{}) string {
	if len(values) == 0 {
		return "string"
	}
	switch values[0].(type) {
	case string:
		return "string"
	case int, int32, int64:
		return "integer"
	case float32, float64:
		return "number"
	default:
		return "string"
	}
}
