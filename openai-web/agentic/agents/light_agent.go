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

func LightAgent(ctx context.Context) (agent.Agent, error) {
	model, err := models.GetModel(ctx)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	getLightStatusTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_lights",
			Description: "Gets a list of lights and their current state.",
		},
		getLightStatus)
	if err != nil {
		log.Fatalf("Failed to create get lights state tool: %v", err)
	}

	changeLightStatusTool, err := functiontool.New(
		functiontool.Config{
			Name:        "change_state",
			Description: "Changes the state of the light and returns all lights.",
		},
		changeLightStatus)
	if err != nil {
		log.Fatalf("Failed to create change lights state tool: %v", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "light_agent",
		Model:       model,
		Description: "Agent to control light's status.",
		Instruction: "You are a useful light assistant. can tall user the status of the lights and can help user control the lights on and off.",
		Tools: []tool.Tool{
			getLightStatusTool, changeLightStatusTool,
		},
	})
}

type getLightStatusInput struct{}

func getLightStatus(_ tool.Context, _ getLightStatusInput) (LightListInfo, error) {
	return sampleLightDatas, nil
}

type changeLightStatusInput struct {
	Id   string `json:"id"`
	IsOn bool   `json:"is_on"`
}

func changeLightStatus(_ tool.Context, input changeLightStatusInput) (LightListInfo, error) {
	result := LightListInfo{Items: []LightInfo{}}
	for _, info := range sampleLightDatas.Items {
		if info.Id == input.Id {
			info.IsOn = input.IsOn
		}
		result.Items = append(result.Items, info)
	}

	return result, nil
}

type LightInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	IsOn bool   `json:"is_on"`
}

type LightListInfo struct {
	Items []LightInfo `json:"items"`
}

var sampleLightDatas LightListInfo = LightListInfo{
	Items: []LightInfo{
		{
			Id:   "1",
			Name: "Table Lamp",
			IsOn: false,
		},
		{
			Id:   "2",
			Name: "Porch light",
			IsOn: false,
		},
		{
			Id:   "3",
			Name: "Chandelier",
			IsOn: true,
		},
	},
}
