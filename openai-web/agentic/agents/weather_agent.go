package agents

import (
	"context"
	"log"

	"github.com/futugyousuzu/go-openai-web/agentic/models"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

func WeatherAgent(ctx context.Context) (agent.Agent, error) {
	model, err := models.GetModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	checkWeatherTool, err := functiontool.New(
		functiontool.Config{
			Name:        "check_weather",
			Description: "Query the weather for a specified location.",
		},
		checkWeather)
	if err != nil {
		log.Fatalf("Failed to create check weather tool: %v", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "weather_agent",
		Model:       model,
		Description: "Agent to answer questions about the weather in a city.",
		Instruction: "Your SOLE purpose is to answer questions about the current weather in a specific city. You MUST refuse to answer any questions unrelated to weather.",
		Tools: []tool.Tool{
			checkWeatherTool,
		},
	})
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
