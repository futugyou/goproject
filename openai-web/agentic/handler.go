package agentic

import (
	"encoding/json"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/events"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type Handler struct {
	returnChan chan<- string
	// ID tracking for event correlation
	threadID   string
	runID      string
	messageID  string
	toolCallID string
	stepID     string
}

func NewHandler(returnChan chan<- string) *Handler {
	return &Handler{
		returnChan: returnChan,
		threadID:   events.GenerateThreadID(),
		runID:      events.GenerateRunID(),
	}
}

func (h *Handler) OnBeforeAgent(ctx agent.CallbackContext) (*genai.Content, error) {
	// Generate new message ID for this LLM interaction
	h.messageID = events.GenerateMessageID()

	// Send run started event
	runStartedEvent := events.NewRunStartedEvent(h.threadID, h.runID)
	if jsonData, err := runStartedEvent.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}

	// Send text message start event
	textStartEvent := events.NewTextMessageStartEvent(h.messageID, events.WithRole("assistant"))
	if jsonData, err := textStartEvent.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}

	return nil, nil
}

func (h *Handler) OnBeforeModel(ctx agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
	h.stepID = events.GenerateStepID()

	// Send step started event with step name
	stepStartedEvent := events.NewStepStartedEvent(h.stepID)
	if jsonData, err := stepStartedEvent.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}

	return nil, nil
}

func (h *Handler) OnBeforeTool(ctx tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
	// Generate tool call ID
	h.toolCallID = events.GenerateToolCallID()

	// Extract tool name from input if possible
	toolName := tool.Name()

	// Send tool call start event
	toolStartEvent := events.NewToolCallStartEvent(h.toolCallID, toolName)
	if h.messageID != "" {
		toolStartEvent = events.NewToolCallStartEvent(h.toolCallID, toolName, events.WithParentMessageID(h.messageID))
	}
	if jsonData, err := toolStartEvent.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}

	input, _ := json.Marshal(args)
	// Send tool arguments
	toolArgsEvent := events.NewToolCallArgsEvent(h.toolCallID, string(input))
	if jsonData, err := toolArgsEvent.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}

	return nil, nil
}

func (h *Handler) OnAfterTool(ctx tool.Context, tool tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
	if h.toolCallID != "" {
		// Send tool call end event
		toolEndEvent := events.NewToolCallEndEvent(h.toolCallID)
		if jsonData, err := toolEndEvent.ToJSON(); err == nil {
			h.returnChan <- string(jsonData)
		}

		output, _ := json.Marshal(result)
		// Send tool call result event
		resultMessageID := events.GenerateMessageID()
		toolResultEvent := events.NewToolCallResultEvent(resultMessageID, h.toolCallID, string(output))
		if jsonData, err := toolResultEvent.ToJSON(); err == nil {
			h.returnChan <- string(jsonData)
		}

		h.toolCallID = ""
	}
	return nil, nil
}

func (h *Handler) OnAfterModel(ctx agent.CallbackContext, llmResponse *model.LLMResponse, llmResponseError error) (*model.LLMResponse, error) {
	if h.stepID != "" {
		stepFinishedEvent := events.NewStepFinishedEvent(h.stepID)
		if jsonData, err := stepFinishedEvent.ToJSON(); err == nil {
			h.returnChan <- string(jsonData)
		}
		h.stepID = ""
	}

	return nil, nil
}

func (h *Handler) OnTextMessaging(part *genai.Part, partial bool) error {
	return nil
}

func (h *Handler) OnAfterAgent(ctx agent.CallbackContext) (*genai.Content, error) {
	runFinishedEvent := events.NewRunFinishedEvent(h.threadID, h.runID)
	if jsonData, err := runFinishedEvent.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}

	return nil, nil
}
