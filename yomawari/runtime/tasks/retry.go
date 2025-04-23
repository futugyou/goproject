package tasks

import (
	"context"
	"time"
)

type RetryPolicy struct {
	MaxRetries int
	Delay      time.Duration
}

//	fut := RunWithRetry(context.Background(), func(ctx context.Context) (int, error) {
//		return 0, errors.New("failed")
//	}, RetryPolicy{
//		MaxRetries: 3,
//		Delay:      500 * time.Millisecond,
//	})
//
// val, err := fut.Result()
// fmt.Println("Result with retry:", val, err)
func RunWithRetry[T any](
	ctx context.Context,
	fn func(context.Context) (T, error),
	policy RetryPolicy,
) Future[T] {
	future, setter := NewFuture[T](ctx)

	go func() {
		defer setter.SetCanceled()

		var result T
		var err error
		for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
			if ctx.Err() != nil {
				return
			}

			result, err = fn(ctx)
			if err == nil {
				setter.SetResult(result)
				return
			}

			time.Sleep(policy.Delay)
		}
		setter.SetError(err)
	}()

	return future
}
