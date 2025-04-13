package shared

import (
	"context"
	"sync"
)

type TaskCompletionSource[T any] struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	responseChan chan T
	err          error
	once         sync.Once
	setCalled    chan struct{}
}

func NewTaskCompletionSource[T any](ctx context.Context, cancelFunc context.CancelFunc) *TaskCompletionSource[T] {
	return &TaskCompletionSource[T]{
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		responseChan: make(chan T, 1),
		setCalled:    make(chan struct{}),
	}
}

func (tcs *TaskCompletionSource[T]) SetResult(result T) {
	tcs.once.Do(func() {
		tcs.responseChan <- result
		close(tcs.setCalled)
	})
}

func (tcs *TaskCompletionSource[T]) SetError(err error) {
	tcs.once.Do(func() {
		tcs.err = err
		close(tcs.responseChan)
		close(tcs.setCalled)
	})
}

func (tcs *TaskCompletionSource[T]) SetCanceled() {
	tcs.Cancel()
	tcs.SetError(context.Canceled)
}

func (tcs *TaskCompletionSource[T]) Cancel() {
	tcs.cancelFunc()
}

func (tcs *TaskCompletionSource[T]) Result() (T, error) {
	var zero T
	select {
	case val, ok := <-tcs.responseChan:
		if !ok {
			return zero, tcs.err
		}
		return val, nil
	case <-tcs.ctx.Done():
		return zero, tcs.ctx.Err()
	}
}

func (tcs *TaskCompletionSource[T]) IsCompleted() bool {
	select {
	case <-tcs.setCalled:
		return true
	default:
		return false
	}
}

func (tcs *TaskCompletionSource[T]) TrySetResult(result T) bool {
	called := false
	tcs.once.Do(func() {
		tcs.responseChan <- result
		close(tcs.setCalled)
		called = true
	})
	return called
}

func (tcs *TaskCompletionSource[T]) TrySetError(err error) bool {
	called := false
	tcs.once.Do(func() {
		tcs.err = err
		close(tcs.responseChan)
		close(tcs.setCalled)
		called = true
	})
	return called
}

func (tcs *TaskCompletionSource[T]) TrySetCanceled() bool {
	tcs.Cancel()
	return tcs.TrySetError(context.Canceled)
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
