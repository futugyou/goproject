package tasks

import (
	"context"
	"errors"
	"time"
)

func TimeoutAfter[T any](f Future[T], duration time.Duration) Future[T] {
	ctx, cancel := context.WithTimeout(context.Background(), duration)

	timedFuture, setter := NewFuture[T](ctx)

	go func() {
		defer cancel()

		select {
		case <-ctx.Done():
			setter.SetError(errors.New("operation timed out"))
		default:
			val, err := f.Result()
			if err != nil {
				setter.SetError(err)
			} else {
				setter.SetResult(val)
			}
		}
	}()

	return timedFuture
}
