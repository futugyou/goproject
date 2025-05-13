package ai

import "context"

type ITextTokenizer interface {
	CountTokens(ctx context.Context, text string) int64
	GetTokens(ctx context.Context, text string) []string
}
