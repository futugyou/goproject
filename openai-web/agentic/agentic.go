package agentic

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/encoding/sse"
	"golang.org/x/sync/errgroup"
)

func ProcessInput(ctx context.Context, w *bufio.Writer, sseWriter *sse.SSEWriter, input string) error {
	resultChan := make(chan string)
	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		for {
			select {
			case result := <-resultChan:
				if result == "" {
					return nil
				}

				// All messages from the handler should now be proper JSON events
				// WriteBytes will format them as SSE frames with "data: " prefix
				if err := sseWriter.WriteBytes(ctx, w, []byte(result)); err != nil {
					return fmt.Errorf("failed to write event: %w", err)
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	g.Go(func() error {
		callLLMErr := CallLLM(groupCtx, input, nil, resultChan)
		close(resultChan)
		return callLLMErr
	})

	return g.Wait()
}
