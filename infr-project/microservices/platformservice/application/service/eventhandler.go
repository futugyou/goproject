package service

import (
	"context"
	"time"

	"github.com/futugyou/domaincore/application"
)

type EventLog struct {
	EventID    string
	PlatformID string
	ProjectID  string
	Token      string
	EventType  string
	CreatedAt  time.Time
}

type EventHandler interface {
	application.EventDispatcher
	Get(ctx context.Context, eventID string) (*EventLog, error)
}
