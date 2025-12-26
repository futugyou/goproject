package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/futugyou/extensions"
	"github.com/futugyousuzu/go-openai-web/agentic"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/encoding/sse"
)

func AguiHandler(w http.ResponseWriter, r *http.Request) {
	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	logger := slog.Default()
	sseWriter := sse.NewSSEWriter().WithLogger(logger)

	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	logCtx := []any{
		"request_id", requestID,
		"route", r.URL.Path,
		"method", r.Method,
	}

	var input AgenticInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		logger.Error("Failed to parse request body", append(logCtx, "error", err)...)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	logger.Info("Tool-based generative UI SSE connection established: %v", logCtx)

	writer := bufio.NewWriter(w)

	ctx := r.Context()

	go func() {
		defer writer.Flush()
		if err := streamAgenticEvents(ctx, writer, sseWriter, &input, logger, logCtx); err != nil {
			logger.Error("Error streaming tool-based generative UI events: %v", err)
		}
	}()

	<-ctx.Done()
}

type AgenticInput struct {
	ThreadID       string           `json:"thread_id"`
	RunID          string           `json:"run_id"`
	State          any              `json:"state"`
	Messages       []map[string]any `json:"messages"`
	Tools          []any            `json:"tools"`
	Context        []any            `json:"context"`
	ForwardedProps any              `json:"forwarded_props"`
}

// streamAgenticEvents implements the tool-based generative UI event sequence
func streamAgenticEvents(reqCtx context.Context, w *bufio.Writer, sseWriter *sse.SSEWriter, input *AgenticInput, logger *slog.Logger, logCtx []any) error {
	// Use IDs from input or generate new ones if not provided
	threadID := input.ThreadID
	if threadID == "" {
		threadID = events.GenerateThreadID()
	}
	runID := input.RunID
	if runID == "" {
		runID = events.GenerateRunID()
	}

	// Create a wrapped context for our operations
	ctx := context.Background()

	// Send RUN_STARTED event
	runStarted := events.NewRunStartedEvent(threadID, runID)
	if err := sseWriter.WriteEvent(ctx, w, runStarted); err != nil {
		return fmt.Errorf("failed to write RUN_STARTED event: %w", err)
	}

	// Check for cancellation
	if err := reqCtx.Err(); err != nil {
		logger.Debug("Client disconnected during RUN_STARTED", append(logCtx, "reason", "context_canceled")...)
		return nil
	}

	// Grab last message from input, will be a user message
	var lastMessage map[string]any
	if len(input.Messages) > 0 {
		lastMessage = input.Messages[len(input.Messages)-1]
	}
	// grab "content" field if it exists
	content, ok := lastMessage["content"].(string)
	if !ok {
		return fmt.Errorf("last message does not have content")
	}

	err := agentic.ProcessInput(ctx, w, sseWriter, content)
	if err != nil {
		return fmt.Errorf("failed to process input: %w", err)

	}

	// Check for cancellation before final event
	if err = reqCtx.Err(); err != nil {
		return fmt.Errorf("client disconnected before RUN_FINISHED: %w", err)
	}

	// Send RUN_FINISHED event
	runFinished := events.NewRunFinishedEvent(threadID, runID)
	if err := sseWriter.WriteEvent(ctx, w, runFinished); err != nil {
		return fmt.Errorf("failed to write RUN_FINISHED event: %w", err)
	}

	return nil
}
