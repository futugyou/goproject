package agents

import (
	"context"
	"errors"

	"google.golang.org/adk/agent"
)

func GetAgentByName(ctx context.Context, name string, handler *Handler) (agent.Agent, error) {
	switch name {
	case "light":
		return WeatherAgent(ctx, handler)
	case "mddoc":
		return MsDocAgent(ctx, handler)
	case "time":
		return TimeAgent(ctx, handler)
	case "weather":
		return WeatherAgent(ctx, handler)
	}

	return nil, errors.New("can not find agent with name: " + name)
}
