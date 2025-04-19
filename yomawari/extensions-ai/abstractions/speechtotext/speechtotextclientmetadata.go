package speechtotext

import "net/url"

type SpeechToTextClientMetadata struct {
	ProviderName   *string
	ProviderUri    *url.URL
	DefaultModelId *string
}

func NewSpeechToTextClientMetadata(providerName string, providerUri *url.URL, defaultModelId string) *SpeechToTextClientMetadata {
	return &SpeechToTextClientMetadata{ProviderName: &providerName, ProviderUri: providerUri, DefaultModelId: &defaultModelId}
}
