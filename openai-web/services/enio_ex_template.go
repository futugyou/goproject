package services

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/futugyousuzu/go-openai-web/models"
)

func (e *EinoService) getChatTemplateNode(ctx context.Context, node models.Node) (prompt.ChatTemplate, error) {
	role := schema.System
	if r, ok := node.Data["role"].(schema.RoleType); ok && len(r) > 0 {
		role = r
	}
	if content, ok := node.Data["content"].(string); ok && len(content) > 0 {
		return prompt.FromMessages(schema.FString, &schema.Message{
			Role:    role,
			Content: content,
		}), nil
	}

	return nil, fmt.Errorf("invalid chat template node: %s", node.ID)
}
