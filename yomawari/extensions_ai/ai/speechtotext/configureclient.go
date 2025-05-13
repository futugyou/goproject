package speechtotext

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/speechtotext"
)

var _ speechtotext.ISpeechToTextClient = (*ConfigureOptionsSpeechToTextClient)(nil)

type ConfigureOptionsSpeechToTextClient struct {
	*speechtotext.DelegatingSpeechToTextClient
	configureOptions func(*speechtotext.SpeechToTextOptions)
	metadata         *speechtotext.SpeechToTextClientMetadata
}

func NewConfigureOptionsSpeechToTextClient(client speechtotext.ISpeechToTextClient, configureOptions func(*speechtotext.SpeechToTextOptions), metadata *speechtotext.SpeechToTextClientMetadata) *ConfigureOptionsSpeechToTextClient {
	return &ConfigureOptionsSpeechToTextClient{
		DelegatingSpeechToTextClient: speechtotext.NewDelegatingSpeechToTextClient(client),
		configureOptions:             configureOptions,
		metadata:                     metadata,
	}
}

// GetStreamingText implements ISpeechToTextClient.
func (d *ConfigureOptionsSpeechToTextClient) GetStreamingText(ctx context.Context, audioSpeechStream io.ReadCloser, options *speechtotext.SpeechToTextOptions) (<-chan speechtotext.SpeechToTextResponse, <-chan error) {
	return d.DelegatingSpeechToTextClient.GetStreamingText(ctx, audioSpeechStream, d.configure(options))
}

// GetStreamingTextWithDataConten implements ISpeechToTextClient.
func (d *ConfigureOptionsSpeechToTextClient) GetStreamingTextWithDataConten(ctx context.Context, audioSpeechContent contents.DataContent, options *speechtotext.SpeechToTextOptions) (<-chan speechtotext.SpeechToTextResponse, <-chan error) {
	return d.DelegatingSpeechToTextClient.GetStreamingTextWithDataConten(ctx, audioSpeechContent, d.configure(options))
}

// GetText implements ISpeechToTextClient.
func (d *ConfigureOptionsSpeechToTextClient) GetText(ctx context.Context, audioSpeechStream io.ReadCloser, options *speechtotext.SpeechToTextOptions) (*speechtotext.SpeechToTextResponse, error) {
	return d.DelegatingSpeechToTextClient.GetText(ctx, audioSpeechStream, d.configure(options))
}

// GetTextWithDataContent implements ISpeechToTextClient.
func (d *ConfigureOptionsSpeechToTextClient) GetTextWithDataContent(ctx context.Context, audioSpeechContent contents.DataContent, options *speechtotext.SpeechToTextOptions) (*speechtotext.SpeechToTextResponse, error) {
	return d.DelegatingSpeechToTextClient.GetTextWithDataContent(ctx, audioSpeechContent, d.configure(options))
}

func (d *ConfigureOptionsSpeechToTextClient) configure(options *speechtotext.SpeechToTextOptions) *speechtotext.SpeechToTextOptions {
	if options == nil {
		options = &speechtotext.SpeechToTextOptions{}
	} else {
		options = options.Clone()
	}

	d.configureOptions(options)
	return options
}
