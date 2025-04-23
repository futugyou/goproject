package tasks

import (
	"context"
)

//	fut := async.RunAsync(context.Background(), func(ctx context.Context) (string, error) {
//		time.Sleep(1 * time.Second)
//		return "hello world", nil
//	})
//
// result, err := fut.Result()
// fmt.Println(result, err)
func RunAsync[T any](ctx context.Context, fn func(context.Context) (T, error)) Future[T] {
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
