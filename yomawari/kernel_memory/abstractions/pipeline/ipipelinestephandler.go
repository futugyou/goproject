package pipeline

import "context"

type IPipelineStepHandler interface {
	Invoke(ctx context.Context, dataPipeline *DataPipeline) (ReturnType, *DataPipeline, error)
	GetStepName() string
	SetStepName(name string)
}

type ReturnType int

const (
	ReturnTypeSuccess ReturnType = iota
	ReturnTypeTransientError
	ReturnTypeFatalError
)
