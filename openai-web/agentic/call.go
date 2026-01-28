package agentic

import (
	"context"
	"fmt"
	"log"

	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/types"
	"github.com/futugyousuzu/go-openai-web/agentic/agents"
	"github.com/futugyousuzu/go-openai-web/agentic/models"
	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func CallLLM(ctx context.Context, input *AgenticInput, returnChan chan<- string) error {
	model, err := models.GetModel(ctx)
	if err != nil {
		return fmt.Errorf("failed to create model: %w", err)
	}

	hander := agents.NewHandler(&input.RunAgentInput, returnChan)

	adkAgent, err := agents.CreateADKAgent(ctx, input.AgentID, model, hander)
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

	r, err := runner.New(runner.Config{AppName: appName, Agent: adkAgent, SessionService: sessionService})
	if err != nil {
		log.Fatalf("FATAL: Failed to create runner: %v", err)
	}

	// Grab last message from input, will be a user message
	var lastMessage types.Message
	if len(input.Messages) > 0 {
		lastMessage = input.Messages[len(input.Messages)-1]
	}

	// grab "content" field if it exists
	content, ok := lastMessage.Content.(string)
	if !ok {
		return fmt.Errorf("last message does not have content")
	}

	userMsg := genai.NewContentFromText(content, genai.RoleUser)
	for event, err := range r.Run(ctx, resp.Session.UserID(), resp.Session.ID(), userMsg, agent.RunConfig{
		StreamingMode: agent.StreamingModeSSE,
	}) {
		if err != nil {
			hander.HandleRunError(err.Error())
		} else {
			if event.Content == nil {
				continue
			}

			for _, part := range event.Content.Parts {
				hander.OnTextMessaging(part, event.Partial)
			}
		}
	}

	return nil
}
