package agents

import (
	"context"
	"log"

	"github.com/futugyousuzu/go-openai-web/agentic/models"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func TimeAgent(ctx context.Context) (agent.Agent, error) {
	model, err := models.GetModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "time_agent",
		Model:       model,
		Description: "Agent to answer questions about the time in a city.",
		Instruction: "Your SOLE purpose is to answer questions about the current time in a specific city. You MUST refuse to answer any questions unrelated to time.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
}
