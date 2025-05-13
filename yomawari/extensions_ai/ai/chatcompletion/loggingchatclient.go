package chatcompletion

import (
	"context"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/ai"
)

type LoggingChatClient struct {
	chatcompletion.DelegatingChatClient
	logger   logger.Logger
	metadata *chatcompletion.ChatClientMetadata
}

func NewLoggingChatClient(innerClient chatcompletion.IChatClient, logger logger.Logger, metadata *chatcompletion.ChatClientMetadata) *LoggingChatClient {
	return &LoggingChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{InnerClient: innerClient},
		logger:               logger,
		metadata:             metadata,
	}
}

func (client *LoggingChatClient) GetResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) (*chatcompletion.ChatResponse, error) {
	if client.logger == nil {
		return client.InnerClient.GetResponse(ctx, chatMessages, options)
	}

	if client.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if client.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			client.logger.Tracef("%s invoked: %s. Options: %s. Metadata: %s.",
				"GetResponse", ai.AsJson(chatMessages), ai.AsJson(options), ai.AsJson(client.metadata))
		} else {
			client.logger.Debugf("%s invoked", "GetResponse")
		}
	}

	response, err := client.InnerClient.GetResponse(ctx, chatMessages, options)

	if err == nil {
		if client.logger.IsOutputLevelEnabled(logger.DebugLevel) {
			if client.logger.IsOutputLevelEnabled(logger.TraceLevel) {
				client.logger.Tracef("%s completed: %s.", "GetResponse", ai.AsJson(response))
			} else {
				client.logger.Debugf("%s completed.", "GetResponse")
			}
		}
	} else {
		client.logger.Errorf("%s failed, err %s.", "GetResponse", err.Error())
	}

	return response, err
}

func (client *LoggingChatClient) GetStreamingResponse(
	ctx context.Context,
	chatMessages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
) <-chan chatcompletion.ChatStreamingResponse {
	if client.logger == nil {
		return client.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
	}

	if client.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if client.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			client.logger.Tracef("%s invoked: %s. Options: %s. Metadata: %s.",
				"GetStreamingResponse", ai.AsJson(chatMessages), ai.AsJson(options), ai.AsJson(client.metadata))
		} else {
			client.logger.Debugf("%s invoked", "GetStreamingResponse")
		}
	}

	responseChan := client.InnerClient.GetStreamingResponse(ctx, chatMessages, options)
	outputChan := make(chan chatcompletion.ChatStreamingResponse)

	go func() {
		defer close(outputChan)

		var responseError error

		for response := range responseChan {
			if response.Err != nil {
				responseError = response.Err
			}

			if responseError != nil {
				client.logger.Errorf("%s failed, err %s.", "GetStreamingResponse", responseError.Error())
				break
			}

			if client.logger.IsOutputLevelEnabled(logger.DebugLevel) {
				if client.logger.IsOutputLevelEnabled(logger.TraceLevel) {
					client.logger.Tracef("%s received update: %s.", "GetStreamingResponse", ai.AsJson(response.Update))
				} else {
					client.logger.Debug("GetStreamingResponse received update.")
				}
			}

			select {
			case outputChan <- response:
			case <-ctx.Done():
				responseError = ctx.Err()
			}
		}
	}()

	return outputChan
}
