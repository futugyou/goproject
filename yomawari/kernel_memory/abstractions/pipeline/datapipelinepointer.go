package pipeline

type DataPipelinePointer struct {
	Index       string   `json:"index"`
	DocumentId  string   `json:"document_id"`
	ExecutionId string   `json:"execution_id"`
	Steps       []string `json:"steps"`
}

func NewDataPipelinePointer(pipeline DataPipeline) *DataPipelinePointer {
	return &DataPipelinePointer{
		Index:       pipeline.Index,
		DocumentId:  pipeline.DocumentId,
		ExecutionId: pipeline.ExecutionId,
		Steps:       pipeline.Steps,
	}
}
