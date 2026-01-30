package agents

import (
	"context"
	"strings"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// Defines lighting information with power specifications
type ComplexLightInfo struct {
	Items []ComplexLightItems `json:"items"`
}

type ComplexLightItems struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	IsOn    bool    `json:"is_on"`
	Wattage float64 `json:"wattage"`
}

// Prompt 1: Based on the current energy consumption, how can I reduce the total power consumption by half? Please thoroughly assess the status of all lights before taking any action.
// Prompt 2: The current power load is very high. Please first check the status of all the lights. If both the 'living room chandelier' and the 'study desk lamp' are on, to prevent a power outage, please turn off only the light with the highest power consumption. If only one light is on, maintain the current state. Please consider your actions carefully before proceeding.
func EnergySavingAgent(ctx context.Context, model model.LLM, handler *Handler) (agent.Agent, error) {
	// Tool 1: Obtain lighting and energy consumption data
	getLightsTool, _ := functiontool.New(functiontool.Config{
		Name:        "get_energy_usage",
		Description: "Obtain the specific status of all lights and the power consumption (in watts) of each light.",
	}, func(_ tool.Context, _ struct{}) (ComplexLightInfo, error) {
		return ComplexLightInfo{
			Items: []ComplexLightItems{
				{"1", "Living room crystal chandelier", true, 200.0},
				{"2", "Study desk lamp", true, 15.0},
				{"3", "Hallway light", true, 40.0},
				{"4", "Kitchen ceiling light", true, 60.0},
			},
		}, nil
	})

	// Tool 2: Control switch
	changeStateTool, _ := functiontool.New(functiontool.Config{
		Name:        "set_light_power",
		Description: "Toggle the light switch state.",
	}, func(_ tool.Context, input changeLightStatusInput) (string, error) {
		return "Success", nil
	})

	// Forcibly trigger Thought's prompt strategy
	promptParts := []string{
		"Role: You are an intelligent power dispatching center with highly developed logical reasoning abilities.",
		"Instruction: When a user makes a complex request, you must activate your internal reasoning logic before taking any action.",
		"[Key Instructions]",
		"  - All your derivations, calculations, and trade-offs in tool selection must be included in your 'thought process(Thought Part)'.",
		"  - The final response to the user (Text Part) should only contain a summary of the operation results, and strictly prohibit the inclusion of any derivation details.",
		"  - Before using the tool, you must verify during the thought process whether the current device state supports the operation.",
		"Thinking Process (Reserved for Thought field):",
		"  - First, analyze the current status of all devices and their known power consumption data.",
		"  - Secondly, the system simulates the execution of user commands to determine whether they will trigger logical conflicts or overloads (such as attempting to turn off a device that is already off).",
		"  - Finally, the most efficient tool chain is planned.",
	}

	config := NewLLMAgentConfig(
		"energy_saver",
		strings.Join(promptParts, "\n\n"),
		"An agent responsible for energy consumption analysis and automated optimization recommendations.",
		model,
		[]tool.Tool{getLightsTool, changeStateTool},
		nil,
		handler,
	)

	return llmagent.New(config)
}
