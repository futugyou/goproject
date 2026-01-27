package agents

import (
	"context"
	"errors"

	"google.golang.org/adk/agent"
)

func GetAgentByName(ctx context.Context, name string) (agent.Agent, error) {
	switch name {
	case "light":
		return WeatherAgent(ctx)
	case "mddoc":
		return MsDocAgent(ctx)
	case "time":
		return TimeAgent(ctx)
	case "weather":
		return WeatherAgent(ctx)
	}

	return nil, errors.New("can not find agent with name: " + name)
}
