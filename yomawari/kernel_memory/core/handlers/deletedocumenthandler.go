package handlers

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type DeleteDocumentHandler struct {
	memoryDbs       []memorystorage.IMemoryDb
	documentStorage documentstorage.IDocumentStorage
	stepName        string
}

func NewDeleteDocumentHandler(stepName string, memoryDbs []memorystorage.IMemoryDb, documentStorage documentstorage.IDocumentStorage) *DeleteDocumentHandler {
	return &DeleteDocumentHandler{
		stepName:        stepName,
		memoryDbs:       memoryDbs,
		documentStorage: documentStorage,
	}
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (d *DeleteDocumentHandler) GetStepName() string {
	return d.stepName
}

// Invoke implements pipeline.IPipelineStepHandler.
func (d *DeleteDocumentHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	var err error
	for _, db := range d.memoryDbs {
		f := models.NewMemoryFilter()
		f.ByDocument(dataPipeline.DocumentId)
		list := db.GetList(ctx, dataPipeline.Index, []models.MemoryFilter{*f}, -1, false)
		go func(db memorystorage.IMemoryDb) {
			select {
			case <-ctx.Done():
				return
			default:
				for record := range list {
					if record.Err != nil {
						err = record.Err
						return
					}
					db.Delete(ctx, dataPipeline.Index, record.Record)
				}
			}
		}(db)
	}

	if err != nil {
		return pipeline.ReturnTypeFatalError, dataPipeline, err
	}

	err = d.documentStorage.EmptyDocumentDirectory(ctx, dataPipeline.Index, dataPipeline.DocumentId)
	if err != nil {
		return pipeline.ReturnTypeFatalError, dataPipeline, err
	}

	return pipeline.ReturnTypeSuccess, dataPipeline, err
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (d *DeleteDocumentHandler) SetStepName(name string) {
	d.stepName = name
}

var _ pipeline.IPipelineStepHandler = (*DeleteDocumentHandler)(nil)
