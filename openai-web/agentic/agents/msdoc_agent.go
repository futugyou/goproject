package agents

import (
	"context"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

func mslearnMCPTransport(ctx context.Context) mcp.Transport {
	return &mcp.StreamableClientTransport{
		Endpoint:   "https://learn.microsoft.com/api/mcp",
		HTTPClient: &http.Client{},
	}
}

func MsDocAgent(ctx context.Context, model model.LLM, handler *Handler) (agent.Agent, error) {
	transport := mslearnMCPTransport(ctx)
	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
	}

	config := NewLLMAgentConfig(
		"msdoc",
		"You help with Microsoft documentation questions. All questions related to Microsoft documentation must first be addressed by using the mcp_tool to obtain the answer before providing a response.",
		"Agent to assist with Microsoft documentation using MCP toolset.",
		model,
		nil,
		[]tool.Toolset{
			mcpToolSet,
		},
		handler,
	)

	return llmagent.New(config)
}
