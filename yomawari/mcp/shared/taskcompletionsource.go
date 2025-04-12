package shared

import (
	"context"
	"sync"
)

type TaskCompletionSource[T any] struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	responseChan chan T
	once         sync.Once
}

func NewTaskCompletionSource[T any](ctx context.Context, cancelFunc context.CancelFunc) *TaskCompletionSource[T] {
	return &TaskCompletionSource[T]{
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		responseChan: make(chan T, 1),
	}
}

func (tcs *TaskCompletionSource[T]) SetResult(result T) {
	tcs.once.Do(func() {
		select {
		case tcs.responseChan <- result:
		default:
		}
	})
}

func (tcs *TaskCompletionSource[T]) Result() (T, error) {
	var zero T
	select {
	case r := <-tcs.responseChan:
		return r, nil
	case <-tcs.ctx.Done():
		return zero, tcs.ctx.Err()
	}
}

func (tcs *TaskCompletionSource[T]) Cancel() {
	tcs.cancelFunc()
}

// cancellationToken.Register(() => { ... })。
func RegisterCancellation(ctx context.Context, fn func()) {
	if ctx.Done() == nil {
		return
	}
	go func() {
		<-ctx.Done()
		fn()
	}()
}

func RegisterCancellationWithArgs[T any](ctx context.Context, arg T, fn func(arg T)) {
	if ctx.Done() == nil {
		return
	}
	go func() {
		<-ctx.Done()
		fn(arg)
	}()
}
