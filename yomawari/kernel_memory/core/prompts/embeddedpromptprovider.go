package prompts

import (
	"context"
	"embed"
	"fmt"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/prompts"
)

//go:embed answer-with-facts.txt summarize.txt
var embeddedFiles embed.FS

type EmbeddedPromptProvider struct {
	fs embed.FS
}

func NewEmbeddedPromptProvider() *EmbeddedPromptProvider {
	return &EmbeddedPromptProvider{fs: embeddedFiles}
}

func (p *EmbeddedPromptProvider) ReadPrompt(ctx context.Context, promptName string) (*string, error) {
	fileName := fmt.Sprintf("%s.txt", promptName)
	content, err := p.fs.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("resource %s not found: %w", fileName, err)
	}

	result := string(content)
	return &result, nil
}

var _ prompts.IPromptProvider = (*EmbeddedPromptProvider)(nil)
