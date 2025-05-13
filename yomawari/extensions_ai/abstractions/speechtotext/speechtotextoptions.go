package speechtotext

type SpeechToTextOptions struct {
	ModelId              *string
	SpeechLanguage       *string
	TextLanguage         *string
	SpeechSampleRate     *int
	AdditionalProperties map[string]interface{}
}

func (s *SpeechToTextOptions) Clone() *SpeechToTextOptions {
	if s == nil {
		return nil
	}
	return &SpeechToTextOptions{
		ModelId:              s.ModelId,
		SpeechLanguage:       s.SpeechLanguage,
		TextLanguage:         s.TextLanguage,
		SpeechSampleRate:     s.SpeechSampleRate,
		AdditionalProperties: s.AdditionalProperties,
	}
}
