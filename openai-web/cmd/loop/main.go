package main

import (
	"context"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/workflowagents/loopagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
)

// go run ./cmd/loop/main.go
func main() {
	ctx := context.Background()

	a, err := agents.WeatherAgent(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	loopAgent, err := loopagent.New(loopagent.Config{
		MaxIterations: 2,
		AgentConfig: agent.Config{
			Name:        "loop_agent",
			Description: "A loop agent that runs sub-agents",
			SubAgents:   []agent.Agent{a},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(loopAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
