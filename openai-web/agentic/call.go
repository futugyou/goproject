package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type UIEvent struct {
	Type      string `json:"type"`
	Content   any    `json:"content,omitempty"`
	IsPartial bool   `json:"is_partial"`
	MessageID string `json:"message_id"`
	Timestamp int64  `json:"timestamp"`
}

func CallLLM(ctx context.Context, input string, tools []any, returnChan chan<- string) error {
	geminiModel, err := gemini.NewModel(ctx, os.Getenv("GEMINI_MODEL_ID"), &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	sendEvent := func(evtType string, content any, partial bool) {
		msg := UIEvent{
			Type:      evtType,
			Content:   content,
			IsPartial: partial,
			Timestamp: time.Now().UnixMilli(),
		}
		if jsonData, err := json.Marshal(msg); err == nil {
			returnChan <- string(jsonData)
		}
	}

	onBeforeModel := func(ctx agent.CallbackContext, req *model.LLMRequest) (*model.LLMResponse, error) {
		sendEvent("status", "Model is thinking...", false)
		return nil, nil
	}

	onBeforeTool := func(ctx tool.Context, tool tool.Tool, args map[string]any) (map[string]any, error) {
		sendEvent("tool_start", map[string]any{"name": tool.Name(), "args": args}, false)
		return nil, nil
	}

	onAfterTool := func(ctx tool.Context, tool tool.Tool, args, result map[string]any, err error) (map[string]any, error) {
		sendEvent("tool_end", map[string]any{"name": tool.Name(), "response": result}, false)
		return nil, nil
	}

	onAfterModel := func(ctx agent.CallbackContext, llmResponse *model.LLMResponse, llmResponseError error) (*model.LLMResponse, error) {
		return nil, nil
	}

	llmCfg := llmagent.Config{
		Name:        "AgUIAgent",
		Instruction: "You are a helpful assistant with tool-calling abilities.",
		Model:       geminiModel,

		BeforeModelCallbacks: []llmagent.BeforeModelCallback{onBeforeModel},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{onBeforeTool},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{onAfterTool},
		AfterModelCallbacks:  []llmagent.AfterModelCallback{onAfterModel},
	}

	adkAgent, err := llmagent.New(llmCfg)
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}

	sessionService := session.InMemoryService()
	userID, appName := "console_user", "console_app"
	resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		log.Fatalf("failed to create the session service: %v", err)
	}

	r, err := runner.New(runner.Config{AppName: "appName", Agent: adkAgent, SessionService: sessionService})
	if err != nil {
		log.Fatalf("FATAL: Failed to create runner: %v", err)
	}

	userMsg := genai.NewContentFromText(input, genai.RoleUser)
	for event, err := range r.Run(ctx, resp.Session.UserID(), resp.Session.ID(), userMsg, agent.RunConfig{
		StreamingMode: agent.StreamingModeSSE,
	}) {
		if err != nil {
			fmt.Printf("\nAGENT_ERROR: %v\n", err)
		} else {
			if event.LLMResponse.Content == nil {
				continue
			}

			for _, part := range event.Content.Parts {
				if len(part.Text) > 0 {
					sendEvent("text", part.Text, event.Partial)
				}
			}
		}
	}

	return nil
}
