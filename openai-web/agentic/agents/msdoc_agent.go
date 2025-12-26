package agents

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

func mslearnMCPTransport(ctx context.Context) mcp.Transport {
	return &mcp.StreamableClientTransport{
		Endpoint:   "https://learn.microsoft.com/api/mcp",
		HTTPClient: &http.Client{},
	}
}

func MsDocAgent(ctx context.Context) (agent.Agent, error) {
	model, err := gemini.NewModel(ctx, os.Getenv("GEMINI_MODEL_ID"), &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	transport := mslearnMCPTransport(ctx)
	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: transport,
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
	}

	return llmagent.New(llmagent.Config{
		Name:        "msdoc_agent",
		Model:       model,
		Description: "You help with Microsoft documentation questions.",
		Instruction: "You help with Microsoft documentation questions. All questions related to Microsoft documentation must first be addressed by using the mcp_tool to obtain the answer before providing a response.",
		Toolsets: []tool.Toolset{
			mcpToolSet,
		},
	})
}
