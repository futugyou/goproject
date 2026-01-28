package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/parallelagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
	"github.com/futugyousuzu/go-openai-web/agentic/models"
)

// AGENT_ERROR: Error 400, Message: Please ensure that function call turn comes immediately after a user turn or after a function response turn.
// go run ./cmd/parallel/main.go
func main() {
	ctx := context.Background()

	model, err := models.GetModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	weatherAgent, err := agents.WeatherAgent(ctx, model, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	docAgent, err := agents.MsDocAgent(ctx, model, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	parallelAgent, err := parallelagent.New(parallelagent.Config{
		AgentConfig: agent.Config{
			Name:        "parallel_agent",
			Description: "A parallel agent that runs sub-agents",
			SubAgents:   []agent.Agent{weatherAgent, docAgent},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(parallelAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
