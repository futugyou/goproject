package handlers

import (
	"fmt"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type FileDetailsWithRecordId struct {
	File     *pipeline.GeneratedFileDetails
	RecordId string
}

func NewFileDetailsWithRecordId(file *pipeline.GeneratedFileDetails, pipe *pipeline.DataPipeline) *FileDetailsWithRecordId {
	recordId := fmt.Sprintf("d=%s//p=%s", pipe.DocumentId, file.Id)
	return &FileDetailsWithRecordId{
		File:     file,
		RecordId: recordId,
	}
}
