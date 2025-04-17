package shared

import (
	"context"

	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

// ClientCapabilities represents the capabilities supported by the client.
type ClientCapabilities struct {
	Experimental         map[string]interface{} `json:"experimental,omitempty"`
	Roots                *RootsCapability       `json:"roots,omitempty"`
	Sampling             *SamplingCapability    `json:"sampling,omitempty"`
	NotificationHandlers map[string]NotificationHandler
}

// RootsCapability represents the roots capability configuration.
type RootsCapability struct {
	ListChanged  *bool                                                                                        `json:"listChanged,omitempty"`
	RootsHandler func(ctx context.Context, req *types.ListRootsRequestParams) (*types.ListRootsResult, error) `json:"-"`
}

// SamplingCapability represents the sampling capability configuration.
type SamplingCapability struct {
	SamplingHandler func(ctx context.Context, req *types.CreateMessageRequestParams, progress ProgressReporter) (*types.CreateMessageResult, error) `json:"-"`
}

// ProgressReporter represents a progress notification mechanism
type ProgressReporter interface {
	Report(value messages.ProgressNotificationValue)
}
