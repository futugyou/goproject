package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	"github.com/google/uuid"
)

var _ pipeline.IPipelineStepHandler = (*TextExtractionHandler)(nil)

type TextExtractionHandler struct {
	orchestrator pipeline.IPipelineOrchestrator
	decoders     []dataformats.IContentDecoder
	webScraper   dataformats.IWebScraper
	stepName     string
}

func NewTextExtractionHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator, decoders []dataformats.IContentDecoder, webScraper dataformats.IWebScraper) *TextExtractionHandler {
	handler := &TextExtractionHandler{
		orchestrator: orchestrator,
		decoders:     decoders,
		webScraper:   webScraper,
		stepName:     stepName,
	}
	return handler
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (s *TextExtractionHandler) GetStepName() string {
	return s.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (s *TextExtractionHandler) SetStepName(name string) {
	s.stepName = name
}

// Invoke implements pipeline.IPipelineStepHandler.
func (s *TextExtractionHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	for i := range dataPipeline.Files {
		uploadedFile := &dataPipeline.Files[i]
		if uploadedFile.AlreadyProcessedBy(s, nil) {
			continue
		}

		var sourceFile = uploadedFile.Name
		var destFile = fmt.Sprintf("%s..extract.txt", uploadedFile.Name)
		var destFile2 = fmt.Sprintf("%s..extract.json", uploadedFile.Name)
		fileContent, err := s.orchestrator.ReadFile(ctx, dataPipeline, sourceFile)
		if err != nil {
			continue
		}
		text := ""
		content := &dataformats.FileContent{MimeType: pipeline.MimeTypes_PlainText}
		skipFile := false
		if len(fileContent) > 0 {
			if uploadedFile.MimeType == pipeline.MimeTypes_WebPageUrl {
				downloadedPage, pageContent, skip := s.DownloadContent(ctx, uploadedFile, fileContent)
				skipFile = skip
				if !skipFile {
					text, content, skipFile, err = s.ExtractText(ctx, downloadedPage, pageContent)
				}
			} else {
				text, content, skipFile, err = s.ExtractText(ctx, uploadedFile, fileContent)
			}
			if err != nil {
				continue
			}
		}

		contentBytes, err := json.Marshal(content)
		if err != nil {
			continue
		}

		if !skipFile {
			s.orchestrator.WriteFile(ctx, dataPipeline, destFile, []byte(text))
			destFileDetails := pipeline.GeneratedFileDetails{
				FileDetailsBase: pipeline.FileDetailsBase{
					Id:           uuid.New().String(),
					Name:         destFile,
					Size:         int64(len(text)),
					MimeType:     content.MimeType,
					ArtifactType: pipeline.ArtifactTypesExtractedText,
					Tags:         dataPipeline.Tags,
					ProcessedBy:  []string{},
					LogEntries:   []pipeline.PipelineLogEntry{},
				},
				ParentId: uploadedFile.Id,
			}
			destFileDetails.MarkProcessedBy(s, nil)
			uploadedFile.GeneratedFiles[destFile] = destFileDetails

			s.orchestrator.WriteFile(ctx, dataPipeline, destFile2, contentBytes)
			destFile2Details := pipeline.GeneratedFileDetails{
				FileDetailsBase: pipeline.FileDetailsBase{
					Id:           uuid.New().String(),
					Name:         destFile2,
					Size:         int64(len(text)),
					MimeType:     content.MimeType,
					ArtifactType: pipeline.ArtifactTypesExtractedContent,
					Tags:         dataPipeline.Tags,
					ProcessedBy:  []string{},
					LogEntries:   []pipeline.PipelineLogEntry{},
				},
				ParentId: uploadedFile.Id,
			}
			destFile2Details.MarkProcessedBy(s, nil)
			uploadedFile.GeneratedFiles[destFile2] = destFile2Details
		}

		uploadedFile.MarkProcessedBy(s, nil)
	}

	return pipeline.ReturnTypeSuccess, dataPipeline, nil
}

// Invoke implements pipeline.IPipelineStepHandler.
func (s *TextExtractionHandler) ExtractText(ctx context.Context, uploadedFile *pipeline.FileDetails, fileContent []byte) (string, *dataformats.FileContent, bool, error) {
	content := &dataformats.FileContent{MimeType: pipeline.MimeTypes_PlainText}
	if uploadedFile == nil || len(uploadedFile.MimeType) == 0 {
		uploadedFile.Log(s, fmt.Sprintf("file MIME type is empty, ignoring the file %s", uploadedFile.Name))
		return "", content, true, nil
	}
	var decoder dataformats.IContentDecoder
	for i := len(s.decoders) - 1; i >= 0; i-- {
		if s.decoders[i].SupportsMimeType(ctx, uploadedFile.MimeType) {
			decoder = s.decoders[i]
			break
		}
	}

	if decoder == nil {
		uploadedFile.Log(s, fmt.Sprintf("file MIME type not supported: %s. Ignoring the file %s", uploadedFile.MimeType, uploadedFile.Name))
		return "", content, true, nil
	} else {
		if con, err := decoder.Decode(ctx, string(fileContent)); err != nil {
			content = con
		}
	}

	textBuilder := &strings.Builder{}
	for _, section := range content.Sections {
		var sectionContent = strings.TrimSpace(section.Content)
		if len(sectionContent) == 0 {
			continue
		}

		textBuilder.WriteString(sectionContent)
		if section.SentencesAreComplete() {
			text.AppendLine(textBuilder)
			text.AppendLine(textBuilder)
		}
	}

	var text = strings.TrimSpace(textBuilder.String())

	return text, content, false, nil
}

func (s *TextExtractionHandler) DownloadContent(ctx context.Context, uploadedFile *pipeline.FileDetails, fileContent []byte) (*pipeline.FileDetails, []byte, bool) {
	url := string(fileContent)
	if len(url) == 0 {
		uploadedFile.Log(s, "The web page URL is empty")
		return uploadedFile, fileContent, true
	}

	urlDownloadResult := s.webScraper.GetContent(ctx, url)
	if !urlDownloadResult.Success && urlDownloadResult.Error != nil {
		uploadedFile.Log(s, fmt.Sprintf("web page download error: %s", urlDownloadResult.Error.Error()))
		return uploadedFile, fileContent, true
	}

	if len(urlDownloadResult.Content) == 0 {
		uploadedFile.Log(s, "the web page has no text content, skipping it")
		return uploadedFile, fileContent, true
	}

	// IMPORTANT: copy by value to avoid editing the source var
	data, err := json.Marshal(uploadedFile)
	if err != nil {
		uploadedFile.Log(s, fmt.Sprintf("web page download error: %s", err.Error()))
		return uploadedFile, fileContent, true
	}
	var fileDetail pipeline.FileDetails
	err = json.Unmarshal(data, &fileDetail)
	if err != nil {
		uploadedFile.Log(s, fmt.Sprintf("web page download error: %s", err.Error()))
		return uploadedFile, fileContent, true
	}

	fileDetail.MimeType = urlDownloadResult.ContentType
	fileDetail.Size = int64(len(urlDownloadResult.Content))

	return &fileDetail, urlDownloadResult.Content, false
}
