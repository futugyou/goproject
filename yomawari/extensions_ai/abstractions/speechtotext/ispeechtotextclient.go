package speechtotext

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

type ISpeechToTextClient interface {
	GetText(ctx context.Context, audioSpeechStream io.ReadCloser, options *SpeechToTextOptions) (*SpeechToTextResponse, error)
	GetTextWithDataContent(ctx context.Context, audioSpeechContent contents.DataContent, options *SpeechToTextOptions) (*SpeechToTextResponse, error)
	GetStreamingText(ctx context.Context, audioSpeechStream io.ReadCloser, options *SpeechToTextOptions) (<-chan SpeechToTextResponse, <-chan error)
	GetStreamingTextWithDataConten(ctx context.Context, audioSpeechContent contents.DataContent, options *SpeechToTextOptions) (<-chan SpeechToTextResponse, <-chan error)
}
