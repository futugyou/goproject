package pipeline

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type ArtifactTypes string

const (
	ArtifactTypesUndefined           ArtifactTypes = "Undefined"
	ArtifactTypesTextPartition       ArtifactTypes = "TextPartition"
	ArtifactTypesExtractedText       ArtifactTypes = "ExtractedText"
	ArtifactTypesTextEmbeddingVector ArtifactTypes = "TextEmbeddingVector"
	ArtifactTypesSyntheticData       ArtifactTypes = "SyntheticData"
	ArtifactTypesExtractedContent    ArtifactTypes = "ExtractedContent"
)

type PipelineLogEntry struct {
	Time   time.Time `json:"t"`
	Source string    `json:"src"`
	Text   string    `json:"text"`
}

type FileDetailsBase struct {
	Id              string                `json:"id"`
	Name            string                `json:"name"`
	Size            int64                 `json:"size"`
	MimeType        string                `json:"mime_type"`
	ArtifactType    ArtifactTypes         `json:"artifact_type"`
	PartitionNumber int64                 `json:"partition_number"`
	SectionNumber   int64                 `json:"section_number"`
	Tags            *models.TagCollection `json:"tags"`
	ProcessedBy     []string              `json:"processed_by"`
	LogEntries      []PipelineLogEntry    `json:"log"`
}

func (f *FileDetailsBase) AlreadyProcessedBy(handler IPipelineStepHandler, subStep *string) bool {
	if f == nil {
		return false
	}
	key := handler.GetStepName()
	if subStep != nil && len(*subStep) > 0 {
		key = key + "/" + *subStep
	}

	for _, v := range f.ProcessedBy {
		if v == key {
			return true
		}
	}

	return false
}

func (f *FileDetailsBase) MarkProcessedBy(handler IPipelineStepHandler, subStep *string) {
	if f == nil {
		return
	}
	key := handler.GetStepName()
	if subStep != nil && len(*subStep) > 0 {
		key = key + "/" + *subStep
	}

	if f.ProcessedBy == nil {
		f.ProcessedBy = make([]string, 0)
	}

	f.ProcessedBy = append(f.ProcessedBy, key)
}

func (f *FileDetailsBase) Log(handler IPipelineStepHandler, text string) {
	if f == nil {
		return
	}

	if f.LogEntries == nil {
		f.LogEntries = make([]PipelineLogEntry, 0)
	}

	f.LogEntries = append(f.LogEntries, PipelineLogEntry{
		Time:   time.Now().UTC(),
		Source: handler.GetStepName(),
		Text:   text,
	})
}

type GeneratedFileDetails struct {
	FileDetailsBase
	ParentId          string `json:"parent_id"`
	SourcePartitionId string `json:"source_partition_id"`
	ContentSHA256     string `json:"content_sha256"`
}

type FileDetails struct {
	FileDetailsBase
	GeneratedFiles map[string]GeneratedFileDetails `json:"generated_files"`
}

func (f *FileDetails) GetPartitionFileName(partitionNumber int64) string {
	if f == nil {
		return ""
	}

	return fmt.Sprintf("%s.partition.%d.txt", f.Name, partitionNumber)
}

func (f *FileDetails) GetHandlerOutputFileName(handler IPipelineStepHandler, index int) string {
	if f == nil || handler == nil {
		return ""
	}

	return fmt.Sprintf("%s.%s.%d.txt", f.Name, handler.GetStepName(), index)
}

type DataPipeline struct {
	Index                     string                `json:"index"`
	DocumentId                string                `json:"document_id"`
	ExecutionId               string                `json:"execution_id"`
	Steps                     []string              `json:"steps"`
	RemainingSteps            []string              `json:"remaining_steps"`
	CompletedSteps            []string              `json:"completed_steps"`
	Tags                      *models.TagCollection `json:"tags"`
	Creation                  time.Time             `json:"creation"`
	LastUpdate                time.Time             `json:"last_update"`
	Files                     []FileDetails         `json:"files"`
	ContextArguments          map[string]any        `json:"args"`
	PreviousExecutionsToPurge []DataPipeline        `json:"previous_executions_to_purge"`
	FilesToUpload             []models.UploadedFile `json:"-"`
	UploadComplete            bool                  `json:"-"`
}

func (dp *DataPipeline) Complete() bool {
	if dp == nil {
		return false
	}

	return len(dp.RemainingSteps) == 0
}

func (dp *DataPipeline) Then(stepName string) *DataPipeline {
	if dp == nil {
		return nil
	}
	dp.Steps = append(dp.Steps, stepName)
	return dp
}

func (dp *DataPipeline) AddUploadFile(name, filename, sourceFile string) *DataPipeline {
	if dp == nil {
		return nil
	}
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return dp
	}
	return dp.AddUploadFileByte(name, filename, content)
}

func (dp *DataPipeline) AddUploadFileByte(name, filename string, content []byte) *DataPipeline {
	if dp == nil {
		return nil
	}
	return dp.AddUploadFileReader(name, filename, bytes.NewReader(content))
}

func (dp *DataPipeline) AddUploadFileReader(name, filename string, content io.Reader) *DataPipeline {
	if dp == nil {
		return nil
	}
	readCloser := io.NopCloser(content)
	return dp.AddUploadFileStream(name, filename, readCloser)
}

func (dp *DataPipeline) AddUploadFileStream(name, filename string, content io.ReadCloser) *DataPipeline {
	if dp == nil {
		return nil
	}
	if seeker, ok := content.(io.ReadSeeker); ok {
		seeker.Seek(0, io.SeekStart)
	} else {
		return dp
	}

	dp.FilesToUpload = append(dp.FilesToUpload, models.UploadedFile{
		FileName:    filename,
		FileContent: content,
	})

	return dp
}

func (dp *DataPipeline) Build() *DataPipeline {
	if dp == nil {
		return nil
	}
	if len(dp.FilesToUpload) > 0 {
		dp.UploadComplete = false
	}

	dp.RemainingSteps = dp.Steps
	dp.Creation = time.Now().UTC()
	dp.LastUpdate = dp.Creation

	if err := dp.Validate(); err != nil {
		fmt.Println("error building data pipeline:", err)
	}

	return dp
}

func (dp *DataPipeline) MoveToNextStep() string {
	if dp == nil || len(dp.RemainingSteps) == 0 {
		return ""
	}

	var stepName = dp.RemainingSteps[0]
	dp.RemainingSteps = dp.RemainingSteps[1:]
	dp.CompletedSteps = append(dp.CompletedSteps, stepName)

	return stepName
}

func (dp *DataPipeline) RollbackToPreviousStep() string {
	if dp == nil || len(dp.CompletedSteps) == 0 {
		return ""
	}

	var stepName = dp.CompletedSteps[len(dp.CompletedSteps)-1]
	dp.CompletedSteps = dp.CompletedSteps[:len(dp.CompletedSteps)-1]
	dp.RemainingSteps = append([]string{stepName}, dp.RemainingSteps...)

	return stepName
}

func (dp *DataPipeline) IsDocumentDeletionPipeline() bool {
	if dp == nil {
		return false
	}

	return len(dp.Steps) == 1 && dp.Steps[0] == constant.PipelineStepsDeleteDocument
}

func (dp *DataPipeline) IsIndexDeletionPipeline() bool {
	if dp == nil {
		return false
	}

	return len(dp.Steps) == 1 && dp.Steps[0] == constant.PipelineStepsDeleteIndex
}

func (dp *DataPipeline) Validate() error {
	if dp == nil {
		return fmt.Errorf("data pipeline is nil")
	}

	if len(dp.DocumentId) == 0 {
		return fmt.Errorf("document id is empty")
	}
	previous := ""
	for _, step := range dp.Steps {
		if len(step) == 0 {
			return fmt.Errorf("step name is empty")
		}

		if previous == step {
			return fmt.Errorf("step name is duplicated: %s", step)
		}

		previous = step
	}
	return nil
}

func (dp *DataPipeline) GetFile(id string) *FileDetails {
	if dp == nil {
		return nil
	}

	for _, v := range dp.Files {
		if v.Id == id {
			return &v
		}
	}

	return nil
}

func (dp *DataPipeline) ToDataPipelineStatus() *models.DataPipelineStatus {
	if dp == nil {
		return nil
	}

	var status = models.DataPipelineStatus{
		Completed:      dp.Complete(),
		Empty:          len(dp.Files) == 0,
		Index:          dp.Index,
		DocumentId:     dp.DocumentId,
		Tags:           dp.Tags,
		Creation:       dp.Creation,
		LastUpdate:     dp.LastUpdate,
		Steps:          dp.Steps,
		RemainingSteps: dp.RemainingSteps,
		CompletedSteps: dp.CompletedSteps,
	}

	return &status
}

func (dp *DataPipeline) GetContext() context.IContext {
	if dp == nil {
		return nil
	}

	return context.NewRequestContext(dp.ContextArguments)
}
