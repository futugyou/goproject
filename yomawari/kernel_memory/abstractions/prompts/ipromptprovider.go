package prompts

import "context"

type IPromptProvider interface {
	ReadPrompt(ctx context.Context, promptName string) (*string, error)
}
