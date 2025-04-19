package speechtotext

import (
	"time"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
)

type SpeechToTextResponseUpdate struct {
	StartTime         *time.Time                     `json:"start_time,omitempty"`
	EndTime           *time.Time                     `json:"end_time,omitempty"`
	ResponseId        *string                        `json:"response_id,omitempty"`
	ModelId           *string                        `json:"model_id,omitempty"`
	Contents          []contents.IAIContent          `json:"contents,omitempty"`
	RawRepresentation interface{}                    `json:"-"`
	Kind              SpeechToTextResponseUpdateKind `json:"kind,omitempty"`
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
