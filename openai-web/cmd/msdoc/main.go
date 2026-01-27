package main

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/console"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
)

// go run ./cmd/msdoc/main.go
func main() {
	ctx := context.Background()

	a, err := agents.MsDocAgent(ctx, nil)
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
