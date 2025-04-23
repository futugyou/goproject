package tasks

import (
	"context"
)

type Future[T any] interface {
	Result() (T, error)
	IsCompleted() bool
	Done() <-chan struct{}
	Context() context.Context
}
