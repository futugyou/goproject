package prompts

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"path"
	"reflect"

	"github.com/futugyou/yomawari/kernel-memory/abstractions/prompts"
)

// EmbeddedPromptProvider implements IPromptProvider using embedded files
type EmbeddedPromptProvider struct {
	fs *embed.FS
}

func NewEmbeddedPromptProvider(fs *embed.FS) *EmbeddedPromptProvider {
	if fs == nil {
		fs = &embed.FS{}
	}
	return &EmbeddedPromptProvider{fs: fs}
}

// ReadPrompt reads an embedded prompt file
func (p *EmbeddedPromptProvider) ReadPrompt(ctx context.Context, promptName string) (*string, error) {
	// Get the package name (similar to namespace in C#)
	pkgPath := path.Dir(reflect.TypeOf(EmbeddedPromptProvider{}).PkgPath())
	if pkgPath == "" {
		return nil, errors.New("unable to determine package path")
	}

	// Construct the file path
	fileName := fmt.Sprintf("%s.txt", promptName)
	filePath := path.Join(pkgPath, fileName)

	// Read the embedded file
	file, err := p.fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("resource %s not found: %w", filePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read resource %s: %w", filePath, err)
	}
	result := string(content)
	return &result, nil
}

var _ prompts.IPromptProvider = (*EmbeddedPromptProvider)(nil)
