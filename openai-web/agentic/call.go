package agentic

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func CallLLM(ctx context.Context, input string, tools []any, returnChan chan<- string) error {
	geminiModel, err := gemini.NewModel(ctx, os.Getenv("GEMINI_MODEL_ID"), &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	hander := NewHandler(returnChan)

	llmCfg := llmagent.Config{
		Name:        "AgUIAgent",
		Instruction: "You are a helpful assistant with tool-calling abilities.",
		Model:       geminiModel,

		BeforeAgentCallbacks: []agent.BeforeAgentCallback{hander.OnBeforeAgent},		
		BeforeModelCallbacks: []llmagent.BeforeModelCallback{hander.OnBeforeModel},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{hander.OnBeforeTool},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{hander.OnAfterTool},
		AfterModelCallbacks:  []llmagent.AfterModelCallback{hander.OnAfterModel},
		AfterAgentCallbacks: []agent.AfterAgentCallback{hander.OnAfterAgent},
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
				hander.OnTextMessaging(part, event.Partial)
			}
		}
	}

	return nil
}
