package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/futugyousuzu/go-openai-web/models"
)

func AddNodesToGraph(ctx context.Context, g *compose.Graph[map[string]any, *schema.Message], nodes []models.Node, embed embedding.Embedder, chatModel model.BaseChatModel) error {
	for _, node := range nodes {
		switch node.Type {
		case "branch":
			n, err := getGraphBranch(ctx, node)
			if err != nil {
				return err
			}
			g.AddBranch(node.ID, n)
		case "template":
			n, err := getChatTemplateNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddChatTemplateNode(node.ID, n)
		case "doc":
			n, err := getDocumentTransformerNode(ctx, node, embed)
			if err != nil {
				return err
			}
			g.AddDocumentTransformerNode(node.ID, n)
		case "embed":
			g.AddEmbeddingNode(node.ID, embed)
		case "graph":
			n, err := getGraphNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddGraphNode(node.ID, n)
		case "indexer":
			n, err := getIndexerNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddIndexerNode(node.ID, n)
		case "lambda":
			n, err := getLambdaNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddLambdaNode(node.ID, n)
		case "loader":
			n, err := getLoaderNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddLoaderNode(node.ID, n)
		case "model":
			g.AddChatModelNode(node.ID, chatModel)
		case "passthrough":
			g.AddPassthroughNode(node.ID)
		case "retriever":
			n, err := getRetrieverNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddRetrieverNode(node.ID, n)
		case "tools":
			n, err := getToolsNode(ctx, node)
			if err != nil {
				return err
			}
			g.AddToolsNode(node.ID, n)
		default:
			// Optionally handle unknown types
			return fmt.Errorf("unknown node type: %s", node.Type)
		}
	}

	return nil
}

func AddEdgesToGraph(g *compose.Graph[map[string]any, *schema.Message], edges []models.Edge) {
	for _, edge := range edges {
		g.AddEdge(edge.Source, edge.Target)
	}

	starts, ends := findStartAndEnd(edges)
	for _, edge := range starts {
		g.AddEdge(compose.START, edge)
	}

	for _, edge := range ends {
		g.AddEdge(edge, compose.END)
	}
}

// Currently there can only be one start and one end
func findStartAndEnd(edges []models.Edge) (starts []string, ends []string) {
	sourceSet := make(map[string]struct{})
	targetSet := make(map[string]struct{})

	for _, e := range edges {
		sourceSet[e.Source] = struct{}{}
		targetSet[e.Target] = struct{}{}
	}

	// start node: appears in sourceSet but not in targetSet
	for s := range sourceSet {
		if _, ok := targetSet[s]; !ok {
			starts = append(starts, s)
		}
	}

	// end node: appears in targetSet but not in sourceSet
	for t := range targetSet {
		if _, ok := sourceSet[t]; !ok {
			ends = append(ends, t)
		}
	}

	return
}