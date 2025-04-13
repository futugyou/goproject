package shared

import (
	"context"
	"log"
	"sync"
	"time"
)

type TaskProcessor[Msg any] struct {
	MsgChan           <-chan Msg
	Handler           func(context.Context, Msg) error
	MaxConcurrency    int
	PerMessageTimeout time.Duration
}

func (p *TaskProcessor[Msg]) Run(ctx context.Context) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, p.MaxConcurrency)

loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping task processor...")
			break loop

		case msg, ok := <-p.MsgChan:
			if !ok {
				log.Println("task channel closed")
				break loop
			}

			sem <- struct{}{}
			wg.Add(1)

			go func(m Msg) {
				defer func() {
					<-sem
					wg.Done()
				}()

				// every task gets its own timeout
				msgCtx, cancel := context.WithTimeout(ctx, p.PerMessageTimeout)
				defer cancel()

				_ = p.Handler(msgCtx, m) // handler need handle error by itself
			}(msg)
		}
	}

	wg.Wait()
	log.Println("all in-flight tasks processed, shutdown complete.")
}
