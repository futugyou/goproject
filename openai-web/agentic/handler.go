package agentic

import (
	"encoding/json"
	"sync"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type Handler struct {
	returnChan chan<- string
	mu         sync.Mutex

	threadID           string
	runID              string
	messageID          string
	toolCallID         string
	toolName           string
	stepID             string
	hasStartedThinking bool
	currentMode        string // "thinking", "messaging", "none"
}

func NewHandler(returnChan chan<- string) *Handler {
	return &Handler{
		returnChan:  returnChan,
		threadID:    events.GenerateThreadID(),
		runID:       events.GenerateRunID(),
		currentMode: "none",
	}
}

func (h *Handler) OnBeforeAgent(ctx agent.CallbackContext) (*genai.Content, error) {
	h.handleEvent(events.NewRunStartedEvent(h.threadID, h.runID))
	return nil, nil
}

func (h *Handler) OnBeforeModel(ctx agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cleanupLifecycle()

	h.stepID = events.GenerateStepID()
	h.messageID = events.GenerateMessageID()

	h.handleEvent(events.NewStepStartedEvent(h.stepID))
	return nil, nil
}

func (h *Handler) OnTextMessaging(part *genai.Part, partial bool) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	isThinking := part.Thought && part.Text != ""
	isNormalText := !part.Thought && part.Text != ""
	isTool := part.FunctionCall != nil || part.ExecutableCode != nil

	if isThinking {
		if h.currentMode != "thinking" {
			h.closeActiveTextContainersInternal()
			h.handleEvent(events.NewThinkingStartEvent().WithTitle("Thinking..."))
			h.hasStartedThinking = true
			h.currentMode = "thinking"
		}
		h.handleEvent(events.NewThinkingTextMessageContentEvent(part.Text))
	}

	if isNormalText {
		if h.currentMode == "thinking" {
			h.handleEvent(events.NewThinkingEndEvent())
			h.hasStartedThinking = false
		}

		h.currentMode = "messaging"
		// use stream(sse)
		// eg. user: hello
		// llm response: Hello! How can, (partial is true)
		// llm response: I help you today?\n (partial is true)
		// llm response: Hello! How can I help you today?\n (partial is false)
		// so, if partial is false, do not send chunk event
		if partial {
			// TextMessage start+content(*)+end === chunks
			h.handleEvent(events.NewTextMessageChunkEvent(&h.messageID, toPtr("assistant"), &part.Text))
		}

		if !partial {
			h.currentMode = "none"
		}
	}

	if isTool {
		h.closeActiveTextContainersInternal()
		evv := events.NewToolCallChunkEvent()
		if len(part.Text) > 0 {
			evv.WithToolCallChunkDelta(part.Text)
		}
		h.handleEvent(evv)
		h.currentMode = "tool"
	}

	return nil
}

func (h *Handler) OnAfterModel(ctx agent.CallbackContext, llmResponse *model.LLMResponse, llmResponseError error) (*model.LLMResponse, error) {
	return nil, nil
}

func (h *Handler) OnAfterAgent(ctx agent.CallbackContext) (*genai.Content, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.cleanupLifecycle()

	h.handleEvent(events.NewRunFinishedEvent(h.threadID, h.runID))
	return nil, nil
}

func (h *Handler) cleanupLifecycle() {
	if h.hasStartedThinking {
		h.handleEvent(events.NewThinkingEndEvent())
		h.hasStartedThinking = false
	}

	if h.stepID != "" {
		h.handleEvent(events.NewStepFinishedEvent(h.stepID))
		h.stepID = ""
	}

	h.currentMode = "none"
}

func (h *Handler) closeActiveTextContainersInternal() {
	if h.currentMode == "thinking" {
		h.handleEvent(events.NewThinkingEndEvent())
		h.hasStartedThinking = false
	}
}

func (h *Handler) OnBeforeTool(ctx tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.toolCallID = events.GenerateToolCallID()
	h.toolName = tool.Name()
	// Use TOOL_CALL_CHUNK (streaming mode) instead of START/ARGS/END
	input, _ := json.Marshal(args)
	chunkEv := events.NewToolCallChunkEvent()
	chunkEv.WithToolCallChunkID(h.toolCallID)
	chunkEv.WithToolCallChunkName(h.toolName)
	if h.messageID != "" {
		chunkEv.WithToolCallChunkParentMessageID(h.messageID)
	}
	chunkEv.WithToolCallChunkDelta(string(input))
	h.handleEvent(chunkEv)
	return nil, nil
}

func (h *Handler) OnAfterTool(ctx tool.Context, tool tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.toolCallID != "" {
		// TOOL_CALL_RESULT must be sent separately (not covered by CHUNK)
		output, _ := json.Marshal(result)
		h.handleEvent(events.NewToolCallResultEvent(events.GenerateMessageID(), h.toolCallID, string(output)))
		h.toolCallID = ""
		h.toolName = ""
	}
	return nil, nil
}

func (h *Handler) handleEvent(ev events.Event) {
	if jsonData, err := ev.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}
}

func toPtr[T any](v T) *T { return &v }
