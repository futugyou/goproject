package functions

import "reflect"

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
