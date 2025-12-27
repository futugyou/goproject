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
	threadID           string
	runID              string
	messageID          string
	toolCallID         string
	toolName           string
	stepID             string
	hasStartedThinking bool
	hasStartedMessage  bool
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
	// Generate new message ID for this LLM interaction
	h.messageID = events.GenerateMessageID()

	// Send run started event
	runStartedEvent := events.NewRunStartedEvent(h.threadID, h.runID)
	h.handleEvent(runStartedEvent)

	return nil, nil
}

func (h *Handler) OnBeforeModel(ctx agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
	h.stepID = events.GenerateStepID()

	// Send step started event with step name
	stepStartedEvent := events.NewStepStartedEvent(h.stepID)
	h.handleEvent(stepStartedEvent)

	return nil, nil
}

func (h *Handler) OnBeforeTool(ctx tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
	// Generate tool call ID
	h.toolCallID = events.GenerateToolCallID()

	// Extract tool name from input if possible
	h.toolName = tool.Name()

	// Send tool call start event
	toolStartEvent := events.NewToolCallStartEvent(h.toolCallID, h.toolName)
	if h.messageID != "" {
		toolStartEvent = events.NewToolCallStartEvent(h.toolCallID, h.toolName, events.WithParentMessageID(h.messageID))
	}
	h.handleEvent(toolStartEvent)

	input, _ := json.Marshal(args)
	// Send tool arguments
	toolArgsEvent := events.NewToolCallArgsEvent(h.toolCallID, string(input))
	h.handleEvent(toolArgsEvent)

	return nil, nil
}

func (h *Handler) OnAfterTool(ctx tool.Context, tool tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
	if h.toolCallID != "" {
		// Send tool call end event
		toolEndEvent := events.NewToolCallEndEvent(h.toolCallID)
		h.handleEvent(toolEndEvent)

		output, _ := json.Marshal(result)
		// Send tool call result event
		resultMessageID := events.GenerateMessageID()
		toolResultEvent := events.NewToolCallResultEvent(resultMessageID, h.toolCallID, string(output))
		h.handleEvent(toolResultEvent)

		h.toolCallID = ""
		h.toolName = ""
	}
	return nil, nil
}

func (h *Handler) OnAfterModel(ctx agent.CallbackContext, llmResponse *model.LLMResponse, llmResponseError error) (*model.LLMResponse, error) {
	if h.hasStartedThinking && h.currentMode == "thinking" {
		ev := events.NewThinkingEndEvent()
		h.handleEvent(ev)
	}

	if h.hasStartedMessage {
		ev := events.NewTextMessageEndEvent(h.messageID)
		h.handleEvent(ev)
	}

	if h.stepID != "" {
		stepFinishedEvent := events.NewStepFinishedEvent(h.stepID)
		h.handleEvent(stepFinishedEvent)
		h.stepID = ""
	}

	h.hasStartedThinking = false
	h.hasStartedMessage = false

	return nil, nil
}

func (h *Handler) handleEvent(ev events.Event) {
	if jsonData, err := ev.ToJSON(); err == nil {
		h.returnChan <- string(jsonData)
	}
}

func (h *Handler) OnTextMessaging(part *genai.Part, partial bool) error {
	isThinking := part.Thought && part.Text != ""
	isNormalText := !part.Thought && part.Text != ""
	isTool := part.FunctionCall != nil || part.ExecutableCode != nil
	var ev events.Event

	if isThinking {
		if !h.hasStartedThinking {
			evv := events.NewThinkingStartEvent()
			title := "Thinking..."
			if len(part.Text) > 0 {
				title = part.Text
			}
			evv.WithTitle(title)
			h.handleEvent(evv)
			h.hasStartedThinking = true
			h.currentMode = "thinking"
		}

		ev = events.NewThinkingTextMessageContentEvent(part.Text)
		h.handleEvent(ev)
	}

	if isNormalText {
		if h.currentMode == "thinking" {
			ev = events.NewThinkingEndEvent()
			h.handleEvent(ev)
			h.currentMode = "none"
		}

		if !h.hasStartedMessage {
			// Send text message start event
			ev = events.NewTextMessageStartEvent(h.messageID, events.WithRole("assistant"))
			h.handleEvent(ev)
			h.hasStartedMessage = true
			h.currentMode = "messaging"
		}

		ev = events.NewTextMessageChunkEvent(toPtr(h.messageID), toPtr("assistant"), toPtr(part.Text))
		h.handleEvent(ev)
	}

	if isTool {
		h.closeActiveTextContainers()
		evv := events.NewToolCallChunkEvent()
		if len(part.Text) > 0 {
			evv.WithToolCallChunkDelta(part.Text)
		}

		h.handleEvent(evv)
	}

	return nil
}

func (h *Handler) closeActiveTextContainers() {
	if h.currentMode == "thinking" {
		ev := events.NewThinkingEndEvent()
		h.handleEvent(ev)
	}
	h.currentMode = "none"
}

func (h *Handler) OnAfterAgent(ctx agent.CallbackContext) (*genai.Content, error) {
	runFinishedEvent := events.NewRunFinishedEvent(h.threadID, h.runID)
	h.handleEvent(runFinishedEvent)

	return nil, nil
}

func toPtr[T any](v T) *T {
	return &v
}
