package tasks

import "context"

func Catch[T any](f Future[T], ctx context.Context, handler func(error) T) Future[T] {
	return RunAsync(ctx, func(ctx context.Context) (T, error) {
		val, err := f.Result()
		if err != nil {
			return handler(err), nil
		}
		return val, nil
	})
}
