package tasks

import (
	"context"
	"log"
	"sync"
	"time"
)

type TaskProcessor[Result any] struct {
	ResultChan        <-chan Result
	Handler           func(context.Context, Result) error
	MaxConcurrency    int
	PerMessageTimeout time.Duration
}

func (p *TaskProcessor[Result]) Run(ctx context.Context) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, p.MaxConcurrency)

loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping task processor...")
			break loop

		case result, ok := <-p.ResultChan:
			if !ok {
				log.Println("task channel closed")
				break loop
			}

			sem <- struct{}{}
			wg.Add(1)

			go func(m Result) {
				defer func() {
					<-sem
					wg.Done()
				}()

				// every task gets its own timeout
				resultCtx, cancel := context.WithTimeout(ctx, p.PerMessageTimeout)
				defer cancel()

				_ = p.Handler(resultCtx, m) // handler need handle error by itself
			}(result)
		}
	}

	wg.Wait()
	log.Println("all in-flight tasks processed, shutdown complete.")
}
