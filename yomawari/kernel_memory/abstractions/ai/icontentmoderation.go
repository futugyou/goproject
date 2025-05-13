package ai

import "context"

type IContentModeration interface {
	IsSafe(ctx context.Context, text string) bool
	IsSafeWiththreshold(ctx context.Context, text string, threshold float64) bool
}
