package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
	"github.com/futugyousuzu/go-openai-web/agentic/models"
)

// go run ./cmd/weather/main.go
func main() {
	ctx := context.Background()

	model, err := models.GetModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	a, err := agents.WeatherAgent(ctx, model, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
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

	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          a,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("failed to create runner: %v", err)
	}

	userMsg := genai.NewContentFromText("What's the weather like in Shanghai?", genai.RoleUser)

	prevText := ""
	for event, err := range r.Run(ctx, userID, resp.Session.ID(), userMsg, agent.RunConfig{
		StreamingMode: agent.StreamingModeSSE,
	}) {
		if err != nil {
			fmt.Printf("\nAGENT_ERROR: %v\n", err)
		} else {
			if event.LLMResponse.Content == nil {
				continue
			}

			text := ""
			for _, p := range event.LLMResponse.Content.Parts {
				text += p.Text
			}

			// In SSE mode, always print partial responses and capture them.
			if !event.IsFinalResponse() {
				fmt.Print(text)
				prevText += text
				continue
			}

			// Only print final response if it doesn't match previously captured text.
			if text != prevText {
				fmt.Print(text)
			}

			prevText = ""
		}
	}
}
