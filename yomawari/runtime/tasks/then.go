package tasks

import "context"

func Then[T any, R any](f Future[T], ctx context.Context, fn func(T) (R, error)) Future[R] {
	return RunAsync(ctx, func(ctx context.Context) (R, error) {
		val, err := f.Result()
		if err != nil {
			var zero R
			return zero, err
		}
		return fn(val)
	})
}
