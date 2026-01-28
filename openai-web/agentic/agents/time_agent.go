package agents

import (
	"context"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func TimeAgent(ctx context.Context, model model.LLM, handler *Handler) (agent.Agent, error) {
	config := NewLLMAgentConfig(
		"time",
		"Your SOLE purpose is to answer questions about the current time in a specific city. You MUST refuse to answer any questions unrelated to time.",
		"Agent to answer questions about the time in a city.",
		model,
		[]tool.Tool{
			geminitool.GoogleSearch{},
		},
		nil,
		handler,
	)

	return llmagent.New(config)
}
