package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/sequentialagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
)

// It's quite interesting. When I asked "What are the current time and weather in Shanghai?", the agent replied,
// "I can only provide weather information, I cannot provide the current time. I'm sorry, I can only provide the current time in Shanghai. I cannot provide weather information."
// This is because I configured the weather agent to only provide weather information and prohibit other functions,
// while the time agent was configured to only provide the time and prohibit other functions.

// go run ./cmd/sequential/main.go
func main() {
	ctx := context.Background()

	weatherAgent, err := agents.WeatherAgent(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	timeAgent, err := agents.TimeAgent(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	sequentialAgent, err := sequentialagent.New(sequentialagent.Config{
		AgentConfig: agent.Config{
			Name:        "sequential_agent",
			Description: "A sequential agent that runs sub-agents",
			SubAgents:   []agent.Agent{weatherAgent, timeAgent},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(sequentialAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
