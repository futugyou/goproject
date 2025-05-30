package protocol

import (
	"context"
)

type NotificationHandler func(ctx context.Context, notification *JsonRpcNotification) error

// ClientCapabilities represents the capabilities supported by the client.
type ClientCapabilities struct {
	Experimental         map[string]interface{} `json:"experimental,omitempty"`
	Roots                *RootsCapability       `json:"roots,omitempty"`
	Sampling             *SamplingCapability    `json:"sampling,omitempty"`
	Elicitation          *ElicitationCapability `json:"elicitation,omitempty"`
	NotificationHandlers map[string]NotificationHandler
}

// RootsCapability represents the roots capability configuration.
type RootsCapability struct {
	ListChanged  *bool                                                                            `json:"listChanged,omitempty"`
	RootsHandler func(ctx context.Context, req *ListRootsRequestParams) (*ListRootsResult, error) `json:"-"`
}

// SamplingCapability represents the sampling capability configuration.
type SamplingCapability struct {
	SamplingHandler func(ctx context.Context, req *CreateMessageRequestParams, progress IProgressReporter) (*CreateMessageResult, error) `json:"-"`
}

// IProgressReporter represents a progress notification mechanism
type IProgressReporter interface {
	Report(value ProgressNotificationValue)
}
