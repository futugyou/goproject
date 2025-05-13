package speechtotext

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

var _ ISpeechToTextClient = (*DelegatingSpeechToTextClient)(nil)

type DelegatingSpeechToTextClient struct {
	InnerClient ISpeechToTextClient
}

func NewDelegatingSpeechToTextClient(client ISpeechToTextClient) *DelegatingSpeechToTextClient {
	return &DelegatingSpeechToTextClient{
		InnerClient: client,
	}
}

// GetStreamingText implements ISpeechToTextClient.
func (d *DelegatingSpeechToTextClient) GetStreamingText(ctx context.Context, audioSpeechStream io.ReadCloser, options *SpeechToTextOptions) (<-chan SpeechToTextResponse, <-chan error) {
	return d.InnerClient.GetStreamingText(ctx, audioSpeechStream, options)
}

// GetStreamingTextWithDataConten implements ISpeechToTextClient.
func (d *DelegatingSpeechToTextClient) GetStreamingTextWithDataConten(ctx context.Context, audioSpeechContent contents.DataContent, options *SpeechToTextOptions) (<-chan SpeechToTextResponse, <-chan error) {
	return d.InnerClient.GetStreamingTextWithDataConten(ctx, audioSpeechContent, options)
}

// GetText implements ISpeechToTextClient.
func (d *DelegatingSpeechToTextClient) GetText(ctx context.Context, audioSpeechStream io.ReadCloser, options *SpeechToTextOptions) (*SpeechToTextResponse, error) {
	return d.InnerClient.GetText(ctx, audioSpeechStream, options)
}

// GetTextWithDataContent implements ISpeechToTextClient.
func (d *DelegatingSpeechToTextClient) GetTextWithDataContent(ctx context.Context, audioSpeechContent contents.DataContent, options *SpeechToTextOptions) (*SpeechToTextResponse, error) {
	return d.InnerClient.GetTextWithDataContent(ctx, audioSpeechContent, options)
}
