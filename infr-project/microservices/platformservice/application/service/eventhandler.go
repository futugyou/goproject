package service

import (
	"context"
	"time"

	coreinfr "github.com/futugyou/domaincore/infrastructure"
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
	coreinfr.EventDispatcher
	Get(ctx context.Context, eventID string) (*EventLog, error)
}
