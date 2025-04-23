package tasks

import (
	"context"
	"sync/atomic"
)

func WhenAny[T any](ctx context.Context, futures ...Future[T]) Future[T] {
	resultFuture, setter := NewFuture[T](ctx)

	var triggered atomic.Bool

	for _, fut := range futures {
		go func(f Future[T]) {
			val, err := f.Result()
			if triggered.CompareAndSwap(false, true) {
				if err != nil {
					setter.SetError(err)
				} else {
					setter.SetResult(val)
				}
			}
		}(fut)
	}

	return resultFuture
}
