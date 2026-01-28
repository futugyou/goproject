package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path"

	"github.com/futugyou/extensions"
	"github.com/futugyousuzu/go-openai-web/agentic"
	verceltool "github.com/futugyousuzu/go-openai-web/vercel"
)

func AguiHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logger := slog.Default()
			sendSSEError(w, logger, fmt.Sprintf("Panic: %v", r))
		}
	}()

	if extensions.Cors(w, r) {
		return
	}

	if !verceltool.AuthForVercel(w, r) {
		return
	}

	_, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	logger := slog.Default()

	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	agentid := r.URL.Query().Get("agentid")
	if len(agentid) == 0 {
		agentid = path.Base(r.URL.Path)
	}

	var input agentic.AgenticInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		sendSSEError(w, logger, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	input.RequestID = requestID
	input.AgentID = agentid

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	w.Header().Set("X-Accel-Buffering", "no")

	logger.Info("Tool-based generative UI SSE connection established", "AgentID", agentid, "RequestID", requestID)

	writer := bufio.NewWriter(w)
	err = streamAgenticEvents(r.Context(), writer, &input, logger)

	writer.Flush()

	if err != nil {
		sendSSEError(w, logger, fmt.Sprintf("Streaming failed: %v", err))
	}
}

// streamAgenticEvents implements the tool-based generative UI event sequence
func streamAgenticEvents(ctx context.Context, w *bufio.Writer, input *agentic.AgenticInput, logger *slog.Logger) error {
	// Check for cancellation
	if err := ctx.Err(); err != nil {
		logger.Debug("Client disconnected during RUN_STARTED", "AgentID", input.AgentID, "RequestID", input.RequestID, "reason", "context_canceled")
		return nil
	}
	return agentic.ProcessInput(ctx, w, input, logger)
}

func sendSSEError(w http.ResponseWriter, logger *slog.Logger, message string) {
	runError := map[string]any{
		"type": "RUN_ERROR",
		"data": map[string]any{
			"message": message,
		},
	}

	jsonData, err := json.Marshal(runError)
	if err != nil {
		logger.Error("Failed to marshal error event", "error", err)
		return
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
	if err != nil {
		logger.Error("Failed to write error to SSE stream", "error", err)
		return
	}

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
