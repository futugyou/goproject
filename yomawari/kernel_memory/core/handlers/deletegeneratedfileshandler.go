package handlers

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type DeleteGeneratedFilesHandler struct {
	documentStorage documentstorage.IDocumentStorage
	stepName        string
}

func NewDeleteGeneratedFilesHandler(stepName string, documentStorage documentstorage.IDocumentStorage) *DeleteGeneratedFilesHandler {
	return &DeleteGeneratedFilesHandler{
		stepName:        stepName,
		documentStorage: documentStorage,
	}
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (d *DeleteGeneratedFilesHandler) GetStepName() string {
	return d.stepName
}

// Invoke implements pipeline.IPipelineStepHandler.
func (d *DeleteGeneratedFilesHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	err := d.documentStorage.EmptyDocumentDirectory(ctx, dataPipeline.Index, dataPipeline.DocumentId)
	if err != nil {
		return pipeline.ReturnTypeFatalError, dataPipeline, err
	}

	return pipeline.ReturnTypeSuccess, dataPipeline, err
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (d *DeleteGeneratedFilesHandler) SetStepName(name string) {
	d.stepName = name
}

var _ pipeline.IPipelineStepHandler = (*DeleteGeneratedFilesHandler)(nil)
