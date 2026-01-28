package agents

import (
	"context"
	"log"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

func WeatherAgent(ctx context.Context, model model.LLM, handler *Handler) (agent.Agent, error) {
	checkWeatherTool, err := functiontool.New(
		functiontool.Config{
			Name:        "check_weather",
			Description: "Query the weather for a specified location.",
		},
		checkWeather)
	if err != nil {
		log.Fatalf("Failed to create check weather tool: %v", err)
	}

	config := NewLLMAgentConfig(
		"weather",
		"Your SOLE purpose is to answer questions about the current weather in a specific city. You MUST refuse to answer any questions unrelated to weather.",
		"Agent to answer questions about the weather in a city.",
		model,
		[]tool.Tool{
			checkWeatherTool,
		},
		nil,
		handler,
	)

	return llmagent.New(config)
}

func checkWeather(ctx tool.Context, input checkWeatherInput) (checkWeatherResult, error) {
	return checkWeatherResult{Location: input.Location, Temperature: 25.4}, nil
}

type checkWeatherInput struct {
	Location string `json:"location"`
}

type checkWeatherResult struct {
	Location    string  `json:"location"`
	Temperature float32 `json:"temperature"`
}
