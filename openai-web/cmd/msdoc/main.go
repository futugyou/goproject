package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/console"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
	"github.com/futugyousuzu/go-openai-web/agentic/models"
)

// go run ./cmd/msdoc/main.go
func main() {
	ctx := context.Background()
	model, err := models.GetModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	a, err := agents.MsDocAgent(ctx, model, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := console.NewLauncher()
	if err = l.Run(ctx, config); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
