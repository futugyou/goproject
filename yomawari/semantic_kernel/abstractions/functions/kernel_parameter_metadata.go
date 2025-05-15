package functions

import "reflect"

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
	// TODO: build Schema
	return InitializedSchema{Inferred: true}
}

type InitializedSchema struct {
	Inferred bool
	Schema   KernelJsonSchema
}
