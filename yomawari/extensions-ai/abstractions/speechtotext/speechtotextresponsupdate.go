package speechtotext

import (
	"time"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
)

type SpeechToTextResponseUpdate struct {
	StartTime            *time.Time                     `json:"start_time,omitempty"`
	EndTime              *time.Time                     `json:"end_time,omitempty"`
	ResponseId           *string                        `json:"response_id,omitempty"`
	ModelId              *string                        `json:"model_id,omitempty"`
	Contents             []contents.IAIContent          `json:"contents,omitempty"`
	AdditionalProperties map[string]interface{}         `json:"additional_properties,omitempty"`
	Kind                 SpeechToTextResponseUpdateKind `json:"kind,omitempty"`
	RawRepresentation    interface{}                    `json:"-"`
}

func NewSpeechToTextResponseUpdateWithContents(contents []contents.IAIContent) *SpeechToTextResponseUpdate {
	return &SpeechToTextResponseUpdate{Contents: contents}
}

func NewSpeechToTextResponseUpdateWithContentText(text string) *SpeechToTextResponseUpdate {
	return &SpeechToTextResponseUpdate{Contents: []contents.IAIContent{contents.NewTextContent(text)}}
}

type SpeechToTextResponseUpdateKind string

const (
	SpeechToTextResponseUpdateKindSessionOpen  = "sessionopen"
	SpeechToTextResponseUpdateKindError        = "error"
	SpeechToTextResponseUpdateKindTextUpdating = "textupdating"
	SpeechToTextResponseUpdateKindTextUpdated  = "textupdated"
	SpeechToTextResponseUpdateKindSessionClose = "sessionclose"
	SpeechToTextResponseUpdateKindUnknown      = "unknown"
)

func NewSpeechToTextResponseUpdateKind(kind string) SpeechToTextResponseUpdateKind {
	switch kind {
	case "sessionopen":
		return SpeechToTextResponseUpdateKindSessionOpen
	case "error":
		return SpeechToTextResponseUpdateKindError
	case "textupdating":
		return SpeechToTextResponseUpdateKindTextUpdating
	case "textupdated":
		return SpeechToTextResponseUpdateKindTextUpdated
	case "sessionclose":
		return SpeechToTextResponseUpdateKindSessionClose
	default:
		return SpeechToTextResponseUpdateKindUnknown
	}
}

func ToSpeechToTextResponse(updates []SpeechToTextResponseUpdate) *SpeechToTextResponse {
	response := &SpeechToTextResponse{
		Contents: []contents.IAIContent{},
	}

	var endTime *time.Time

	for _, update := range updates {
		if response.StartTime == nil {
			response.StartTime = update.StartTime
		}

		if update.EndTime != nil {
			endTime = update.EndTime
		}
		contents, responseId, modelId, additionalProperties := processUpdate(update, response.Contents)
		response.Contents = contents
		response.ResponseId = responseId
		response.ModelId = modelId
		response.AdditionalProperties = additionalProperties
		response.EndTime = endTime
	}
	response.Contents = chatcompletion.CoalesceTextContent(response.Contents)
	return response
}

func processUpdate(update SpeechToTextResponseUpdate, conts []contents.IAIContent) (contents []contents.IAIContent, responseId *string, modelId *string, additionalProperties map[string]interface{}) {
	contents = conts
	if update.ResponseId != nil {
		responseId = update.ResponseId
	}

	if update.ModelId != nil {
		modelId = update.ModelId
	}

	additionalProperties = map[string]interface{}{}
	if update.Contents != nil {
		contents = append(contents, update.Contents...)
	}

	if update.AdditionalProperties != nil {
		for key, value := range update.AdditionalProperties {
			additionalProperties[key] = value
		}
	}
	return
}
