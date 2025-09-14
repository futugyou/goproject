package graph

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/futugyousuzu/go-openai-web/models"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func getToolsNode(ctx context.Context, node models.Node) (*compose.ToolsNode, error) {
	tools := []tool.BaseTool{}
	if googletool, ok := node.Data["googlesearch"].(string); ok && len(googletool) > 0 {
		googleTool, err := googlesearch.NewTool(ctx, &googlesearch.Config{
			APIKey:         os.Getenv("GOOGLE_API_KEY"),
			SearchEngineID: os.Getenv("GOOGLE_SEARCH_ENGINE_ID"),
		})
		if err != nil {
			return nil, err
		}
		tools = append(tools, googleTool)
	}

	if googletool, ok := node.Data["duckduckgo"].(string); ok && len(googletool) > 0 {
		searchTool, err := duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{})
		if err != nil {
			return nil, err
		}
		tools = append(tools, searchTool)
	}

	if mcptoolurl, ok := node.Data["mcptoolurl"].(string); ok && len(mcptoolurl) > 0 {
		cli, err := client.NewSSEMCPClient(mcptoolurl)
		if err != nil {
			return nil, err
		}
		err = cli.Start(ctx)
		if err != nil {
			return nil, err
		}

		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "enio-client",
			Version: "1.0.0",
		}

		_, err = cli.Initialize(ctx, initRequest)
		if err != nil {
			return nil, err
		}

		mcpTools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
		if err != nil {
			return nil, err
		}
		tools = append(tools, mcpTools...)
	}

	if len(tools) > 0 {
		return compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
			Tools: tools,
		})
	}

	return nil, fmt.Errorf("invalid tool node: %s", node.ID)
}
