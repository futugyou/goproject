package service

import (
	"context"
	"time"
)

type EventLog struct {
	EventID    string
	PlatformID string
	ProjectID  string
	Token      string
	EventType  string
	CreatedAt  time.Time
}

type EventLogService interface {
	Create(ctx context.Context, event EventLog) error
	Get(ctx context.Context, eventID string) (*EventLog, error)
}
