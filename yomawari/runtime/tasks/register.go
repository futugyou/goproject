package tasks

import "context"

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
