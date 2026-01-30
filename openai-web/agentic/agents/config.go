package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func NewLLMAgentConfig(name string, instruction string, description string, geminiModel model.LLM, tools []tool.Tool, toolsets []tool.Toolset, handler *Handler) llmagent.Config {
	llmCfg := llmagent.Config{
		Name:        name,
		Instruction: instruction,
		Model:       geminiModel,
		Description: description,
		GenerateContentConfig: &genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				IncludeThoughts: true,
				ThinkingLevel:   genai.ThinkingLevelHigh,
			},
		},
	}

	if handler != nil {
		llmCfg.BeforeAgentCallbacks = []agent.BeforeAgentCallback{handler.OnBeforeAgent}
		llmCfg.BeforeModelCallbacks = []llmagent.BeforeModelCallback{handler.OnBeforeModel}
		llmCfg.BeforeToolCallbacks = []llmagent.BeforeToolCallback{handler.OnBeforeTool}
		llmCfg.AfterToolCallbacks = []llmagent.AfterToolCallback{handler.OnAfterTool}
		llmCfg.AfterModelCallbacks = []llmagent.AfterModelCallback{handler.OnAfterModel}
		llmCfg.AfterAgentCallbacks = []agent.AfterAgentCallback{handler.OnAfterAgent}
	}

	if len(tools) > 0 {
		llmCfg.Tools = tools
	}

	if len(toolsets) > 0 {
		llmCfg.Toolsets = toolsets
	}

	return llmCfg
}
