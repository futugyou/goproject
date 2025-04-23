package tasks

import (
	"context"
	"sync"
)

type AsyncQueue[T any] struct {
	maxConcurrency int
	sem            chan struct{}
	mu             sync.Mutex
}

func NewAsyncQueue[T any](maxConcurrency int) *AsyncQueue[T] {
	return &AsyncQueue[T]{
		maxConcurrency: maxConcurrency,
		sem:            make(chan struct{}, maxConcurrency),
	}
}

func (q *AsyncQueue[T]) RunAsync(ctx context.Context, fn func(context.Context) (T, error)) Future[T] {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.sem <- struct{}{}
	defer func() {
		<-q.sem
	}()

	future, setter := NewFuture[T](ctx)

	go func() {
		val, err := fn(ctx)
		if err != nil {
			setter.SetError(err)
		} else {
			setter.SetResult(val)
		}
	}()

	return future
}
