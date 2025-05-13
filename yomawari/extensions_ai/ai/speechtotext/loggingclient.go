package speechtotext

import (
	"context"
	"io"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/speechtotext"
	"github.com/futugyou/yomawari/extensions_ai/ai"
)

var _ speechtotext.ISpeechToTextClient = (*LoggingSpeechToTextClient)(nil)

type LoggingSpeechToTextClient struct {
	*speechtotext.DelegatingSpeechToTextClient
	logger   logger.Logger
	metadata *speechtotext.SpeechToTextClientMetadata
}

func NewLoggingSpeechToTextClient(
	innerClient speechtotext.ISpeechToTextClient,
	logger logger.Logger,
	metadata *speechtotext.SpeechToTextClientMetadata,
) *LoggingSpeechToTextClient {
	return &LoggingSpeechToTextClient{
		DelegatingSpeechToTextClient: speechtotext.NewDelegatingSpeechToTextClient(innerClient),
		logger:                       logger,
		metadata:                     metadata,
	}
}

// GetStreamingText implements ISpeechToTextClient.
func (d *LoggingSpeechToTextClient) GetStreamingText(ctx context.Context, audioSpeechStream io.ReadCloser, options *speechtotext.SpeechToTextOptions) (<-chan speechtotext.SpeechToTextResponse, <-chan error) {
	if d.logger == nil {
		return d.DelegatingSpeechToTextClient.GetStreamingText(ctx, audioSpeechStream, options)
	}

	if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			d.logger.Tracef("%s invoked: %s. Options: %s. Metadata: %s.",
				"GetStreamingText", ai.AsJson(audioSpeechStream), ai.AsJson(options), ai.AsJson(d.metadata))
		} else {
			d.logger.Debugf("%s invoked", "GetStreamingText")
		}
	}

	responseChan, errChan := d.DelegatingSpeechToTextClient.GetStreamingText(ctx, audioSpeechStream, options)

	outputChan := make(chan speechtotext.SpeechToTextResponse)
	outputErrorChan := make(chan error, 1)

	go func() {
		defer close(outputChan)
		defer close(outputErrorChan)

		for {
			select {
			case <-ctx.Done():
				outputErrorChan <- ctx.Err()
				return

			case err, ok := <-errChan:
				if ok && err != nil {
					d.logger.Errorf("%s failed, err %s.", "GetStreamingText", err.Error())
					outputErrorChan <- err
				}
				return

			case response, ok := <-responseChan:
				if !ok {
					return
				}

				if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
					if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
						d.logger.Tracef("%s received update: %s.", "GetStreamingText", ai.AsJson(response))
					} else {
						d.logger.Debug("GetStreamingText received update.")
					}
				}

				select {
				case outputChan <- response:
					// 写入成功
				case <-ctx.Done():
					outputErrorChan <- ctx.Err()
					return
				}
			}
		}
	}()

	return outputChan, outputErrorChan
}

// GetStreamingTextWithDataConten implements ISpeechToTextClient.
func (d *LoggingSpeechToTextClient) GetStreamingTextWithDataConten(ctx context.Context, audioSpeechContent contents.DataContent, options *speechtotext.SpeechToTextOptions) (<-chan speechtotext.SpeechToTextResponse, <-chan error) {
	if d.logger == nil {
		return d.DelegatingSpeechToTextClient.GetStreamingTextWithDataConten(ctx, audioSpeechContent, options)
	}

	if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			d.logger.Tracef("%s invoked: %s. Options: %s. Metadata: %s.",
				"GetStreamingTextWithDataConten", ai.AsJson(audioSpeechContent), ai.AsJson(options), ai.AsJson(d.metadata))
		} else {
			d.logger.Debugf("%s invoked", "GetStreamingTextWithDataConten")
		}
	}

	responseChan, errChan := d.DelegatingSpeechToTextClient.GetStreamingTextWithDataConten(ctx, audioSpeechContent, options)

	outputChan := make(chan speechtotext.SpeechToTextResponse)
	outputErrorChan := make(chan error, 1)

	go func() {
		defer close(outputChan)
		defer close(outputErrorChan)

		for {
			select {
			case <-ctx.Done():
				outputErrorChan <- ctx.Err()
				return

			case err, ok := <-errChan:
				if ok && err != nil {
					d.logger.Errorf("%s failed, err %s.", "GetStreamingText", err.Error())
					outputErrorChan <- err
				}
				return

			case response, ok := <-responseChan:
				if !ok {
					return
				}

				if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
					if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
						d.logger.Tracef("%s received update: %s.", "GetStreamingText", ai.AsJson(response))
					} else {
						d.logger.Debug("GetStreamingText received update.")
					}
				}

				select {
				case outputChan <- response:
					// 写入成功
				case <-ctx.Done():
					outputErrorChan <- ctx.Err()
					return
				}
			}
		}
	}()

	return outputChan, outputErrorChan
}

// GetText implements ISpeechToTextClient.
func (d *LoggingSpeechToTextClient) GetText(ctx context.Context, audioSpeechStream io.ReadCloser, options *speechtotext.SpeechToTextOptions) (*speechtotext.SpeechToTextResponse, error) {
	if d.logger == nil {
		return d.DelegatingSpeechToTextClient.GetText(ctx, audioSpeechStream, options)
	}

	if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			d.logger.Tracef("%s invoked: Options: %s. Metadata: %s.", "GetText", ai.AsJson(options), ai.AsJson(d.metadata))
		} else {
			d.logger.Debugf("%s invoked", "GetText")
		}
	}

	result, err := d.DelegatingSpeechToTextClient.GetText(ctx, audioSpeechStream, options)

	if err != nil {
		d.logger.Errorf("%s failed, err %s.", "GetText", err.Error())
	} else {
		if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
			if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
				d.logger.Tracef("%s completed: %s.", "GetText", ai.AsJson(result))
			} else {
				d.logger.Debugf("%s completed.", "GetText")
			}
		}
	}

	return result, err
}

// GetTextWithDataContent implements ISpeechToTextClient.
func (d *LoggingSpeechToTextClient) GetTextWithDataContent(ctx context.Context, audioSpeechContent contents.DataContent, options *speechtotext.SpeechToTextOptions) (*speechtotext.SpeechToTextResponse, error) {
	if d.logger == nil {
		return d.DelegatingSpeechToTextClient.GetTextWithDataContent(ctx, audioSpeechContent, options)
	}
	if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			d.logger.Tracef("%s invoked: Options: %s. Metadata: %s.Content: %s", "GetTextWithDataContent", ai.AsJson(options), ai.AsJson(d.metadata), ai.AsJson(audioSpeechContent))
		} else {
			d.logger.Debugf("%s invoked", "GetTextWithDataContent")
		}
	}

	result, err := d.DelegatingSpeechToTextClient.GetTextWithDataContent(ctx, audioSpeechContent, options)

	if err != nil {
		d.logger.Errorf("%s failed, err %s.", "GetTextWithDataContent", err.Error())
	} else {
		if d.logger.IsOutputLevelEnabled(logger.DebugLevel) {
			if d.logger.IsOutputLevelEnabled(logger.TraceLevel) {
				d.logger.Tracef("%s completed: %s.", "GetTextWithDataContent", ai.AsJson(result))
			} else {
				d.logger.Debugf("%s completed.", "GetTextWithDataContent")
			}
		}
	}

	return result, err
}
