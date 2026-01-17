package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/server/adkrest"
	"google.golang.org/adk/session"

	"github.com/futugyousuzu/go-openai-web/agentic/agents"
)

// go run ./cmd/rest/main.go
func main() {
	ctx := context.Background()

	a, err := agents.WeatherAgent(ctx)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader:    agent.NewSingleLoader(a),
		SessionService: session.InMemoryService(),
	}

	apiHandler := adkrest.NewHandler(config, 120*time.Second)
	
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
