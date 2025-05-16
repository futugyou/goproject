package functions

import "reflect"

type KernelReturnParameterMetadata struct {
	Name          string
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
		Name:          meta.Name,
	}
}

func (meta KernelReturnParameterMetadata) GetSchema() KernelJsonSchema {
	return InferSchema(meta.ParameterType, nil, meta.Name, meta.Description).Schema
}
