package handlers

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type DeleteIndexHandler struct {
	memoryDbs       []memorystorage.IMemoryDb
	documentStorage documentstorage.IDocumentStorage
	stepName        string
}

func NewDeleteIndexHandler(stepName string, memoryDbs []memorystorage.IMemoryDb, documentStorage documentstorage.IDocumentStorage) *DeleteIndexHandler {
	return &DeleteIndexHandler{
		stepName:        stepName,
		memoryDbs:       memoryDbs,
		documentStorage: documentStorage,
	}
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (d *DeleteIndexHandler) GetStepName() string {
	return d.stepName
}

// Invoke implements pipeline.IPipelineStepHandler.
func (d *DeleteIndexHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	for _, db := range d.memoryDbs {
		db.DeleteIndex(ctx, dataPipeline.Index)
	}

	err := d.documentStorage.DeleteIndexDirectory(ctx, dataPipeline.Index)
	if err != nil {
		return pipeline.ReturnTypeFatalError, dataPipeline, err
	}

	return pipeline.ReturnTypeSuccess, dataPipeline, err
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (d *DeleteIndexHandler) SetStepName(name string) {
	d.stepName = name
}

var _ pipeline.IPipelineStepHandler = (*DeleteIndexHandler)(nil)
