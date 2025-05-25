package abstractions

import (
	"encoding/json"
	"reflect"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/utilities"
)

type KernelParameterMetadata struct {
	Name          string
	Description   string
	DefaultValue  any
	IsRequired    bool
	ParameterType reflect.Type
	Schema        KernelJsonSchema
}

func NewKernelParameterMetadata(name string) *KernelParameterMetadata {
	return &KernelParameterMetadata{
		Name: name,
	}
}

func InferSchema(parameterType reflect.Type, defaultValue any, description string) InitializedSchema {
	var schema *KernelJsonSchema

	if parameterType != nil {
		invalidAsGeneric := parameterType.Kind() == reflect.Func ||
			parameterType.Kind() == reflect.Ptr && parameterType.Elem().Kind() == reflect.UnsafePointer ||
			parameterType.Kind() == reflect.UnsafePointer ||
			parameterType.Kind() == reflect.Invalid

		if !invalidAsGeneric {
			stringDefault := convertToString(defaultValue)
			if stringDefault != "" {
				if description != "" {
					description += " "
				}
				description += "(default value: " + stringDefault + ")"
			}

			var err error

			schema, err = buildSchema(parameterType, description, defaultValue)

			if err != nil {
				schema = nil
			}
		}
	}

	if schema != nil {
		return InitializedSchema{Inferred: true, Schema: *schema}
	}

	return InitializedSchema{Inferred: true}
}

func buildSchema(parameterType reflect.Type, description string, defaultVal interface{}) (*KernelJsonSchema, error) {
	jsonSchema, err := utilities.CreateJsonSchema(parameterType, description, defaultVal, utilities.DefaultAIJsonSchemaCreateOptions)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(jsonSchema)
	if err != nil {
		return nil, err
	}

	return KernelJsonSchemaParseFromBytes(data)
}

func convertToString(defaultValue any) string {
	if condition, ok := defaultValue.(string); ok {
		return condition
	}

	return ""
}

type InitializedSchema struct {
	Inferred bool
	Schema   KernelJsonSchema
}

type KernelReturnParameterMetadata struct {
	Description   string
	ParameterType reflect.Type
	schema        KernelJsonSchema
}

func NewKernelReturnParameterMetadata() *KernelReturnParameterMetadata {
	return &KernelReturnParameterMetadata{}
}

func KernelReturnParameterMetadataClone(meta KernelReturnParameterMetadata) *KernelReturnParameterMetadata {
	return &KernelReturnParameterMetadata{
		Description:   meta.Description,
		ParameterType: meta.ParameterType,
		schema:        meta.schema,
	}
}

func (meta KernelReturnParameterMetadata) GetSchema() KernelJsonSchema {
	return InferSchema(meta.ParameterType, nil, meta.Description).Schema
}

type KernelFunctionMetadata struct {
	Name                 string
	PluginName           string
	Description          string
	Parameters           []KernelParameterMetadata
	ReturnParameter      KernelReturnParameterMetadata
	AdditionalProperties map[string]interface{}
}

func NewKernelFunctionMetadata(name string) *KernelFunctionMetadata {
	return &KernelFunctionMetadata{
		Name: name,
	}
}

func KernelFunctionMetadataClone(meta KernelFunctionMetadata) *KernelFunctionMetadata {
	return &KernelFunctionMetadata{
		Name:                 meta.Name,
		PluginName:           meta.PluginName,
		Description:          meta.Description,
		Parameters:           meta.Parameters,
		ReturnParameter:      meta.ReturnParameter,
		AdditionalProperties: meta.AdditionalProperties,
	}
}
