package agents

import (
	"context"
	"errors"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
)

func CreateADKAgent(ctx context.Context, name string, model model.LLM, handler *Handler) (agent.Agent, error) {
	switch name {
	case "light":
		return LightAgent(ctx, model, handler)
	case "mddoc":
		return MsDocAgent(ctx, model, handler)
	case "time":
		return TimeAgent(ctx, model, handler)
	case "weather":
		return WeatherAgent(ctx, model, handler)
	case "energy_saver":
		return EnergySavingAgent(ctx, model, handler)
	}

	return nil, errors.New("can not find agent with name: " + name)
}
