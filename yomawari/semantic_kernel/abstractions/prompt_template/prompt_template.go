package prompt_template

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/functions"
)

type IPromptTemplate interface {
	Render(ctx context.Context, kernel abstractions.Kernel, arguments *functions.KernelArguments) (*string, error)
}
