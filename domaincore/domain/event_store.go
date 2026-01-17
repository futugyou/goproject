package domain

import (
	"context"
 
)

type EventStore[E DomainEvent] interface {
	Save(ctx context.Context, events []E) error
	Load(ctx context.Context, id string) ([]E, error)
	LoadGreaterthanVersion(ctx context.Context, id string, version int) ([]E, error)
}
