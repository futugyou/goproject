package tasks

import (
	"context"
	"sync"
)

func WhenAll[T any](ctx context.Context, futures ...Future[T]) Future[[]T] {
	allFuture, setter := NewFuture[[]T](ctx)

	go func() {
		var wg sync.WaitGroup
		results := make([]T, len(futures))
		errorsOccurred := make([]error, len(futures))
		anyFailed := false

		for i, fut := range futures {
			wg.Add(1)
			go func(i int, f Future[T]) {
				defer wg.Done()
				val, err := f.Result()
				if err != nil {
					errorsOccurred[i] = err
					anyFailed = true
				} else {
					results[i] = val
				}
			}(i, fut)
		}

		wg.Wait()

		if anyFailed {
			for _, e := range errorsOccurred {
				if e != nil {
					setter.SetError(e)
					return
				}
			}
		} else {
			setter.SetResult(results)
		}
	}()

	return allFuture
}
