package agents

import (
	"github.com/ag-ui-protocol/ag-ui/sdks/community/go/pkg/core/types"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

func NewLLMAgentConfig(name string, instruction string, description string, geminiModel model.LLM, tools []tool.Tool, toolsets []tool.Toolset, input *types.RunAgentInput, returnChan chan<- string) llmagent.Config {
	hander := NewHandler(input, returnChan)

	llmCfg := llmagent.Config{
		Name:                 name,
		Instruction:          instruction,
		Model:                geminiModel,
		Description:          description,
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{hander.OnBeforeAgent},
		BeforeModelCallbacks: []llmagent.BeforeModelCallback{hander.OnBeforeModel},
		BeforeToolCallbacks:  []llmagent.BeforeToolCallback{hander.OnBeforeTool},
		AfterToolCallbacks:   []llmagent.AfterToolCallback{hander.OnAfterTool},
		AfterModelCallbacks:  []llmagent.AfterModelCallback{hander.OnAfterModel},
		AfterAgentCallbacks:  []agent.AfterAgentCallback{hander.OnAfterAgent},
	}

	if len(toolsets) > 0 {
		llmCfg.Tools = tools
	}

	if len(tools) > 0 {
		llmCfg.Toolsets = toolsets
	}

	return llmCfg
}
