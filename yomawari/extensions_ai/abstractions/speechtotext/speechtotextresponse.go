package speechtotext

import (
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

type SpeechToTextResponse struct {
	StartTime            *time.Time             `json:"start_time,omitempty"`
	EndTime              *time.Time             `json:"end_time,omitempty"`
	ResponseId           *string                `json:"response_id,omitempty"`
	ModelId              *string                `json:"model_id,omitempty"`
	Contents             []contents.IAIContent  `json:"contents,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additional_properties,omitempty"`
	RawRepresentation    interface{}            `json:"-"`
}

func (s *SpeechToTextResponse) Text() string {
	if s == nil || len(s.Contents) == 0 {
		return ""
	}

	return contents.ConcatTextContents(s.Contents)
}

func (s *SpeechToTextResponse) ToSpeechToTextResponseUpdates() []SpeechToTextResponseUpdate {
	if s == nil || len(s.Contents) == 0 {
		return []SpeechToTextResponseUpdate{}
	}

	sp := SpeechToTextResponseUpdate{
		StartTime:            s.StartTime,
		EndTime:              s.EndTime,
		ResponseId:           s.ResponseId,
		ModelId:              s.ModelId,
		Contents:             s.Contents,
		RawRepresentation:    s.RawRepresentation,
		Kind:                 SpeechToTextResponseUpdateKindTextUpdated,
		AdditionalProperties: s.AdditionalProperties,
	}
	return []SpeechToTextResponseUpdate{sp}
}
