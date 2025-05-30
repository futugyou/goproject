package server

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

var _ chatcompletion.IChatClient = (*SamplingChatClient)(nil)

type SamplingChatClient struct {
	server IMcpServer
}

func NewSamplingChatClient(server IMcpServer) *SamplingChatClient {
	return &SamplingChatClient{server: server}
}

// GetResponse implements chatcompletion.IChatClient.
func (s *SamplingChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	return s.server.SampleWithChatMessage(ctx, chatMessages, options)
}

// GetStreamingResponse implements chatcompletion.IChatClient.
func (s *SamplingChatClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	response, err := s.GetResponse(ctx, chatMessages, options)
	streamResp := make(chan chatcompletion.ChatStreamingResponse)
	if err != nil {
		streamResp <- chatcompletion.ChatStreamingResponse{
			Update: nil,
			Err:    err,
		}
		close(streamResp)
		return streamResp
	}
	updates := response.ToChatResponseUpdates()
	go func() {
		defer close(streamResp)
		for _, item := range updates {
			streamResp <- chatcompletion.ChatStreamingResponse{
				Update: &item,
				Err:    nil,
			}
		}
	}()
	return streamResp
}
