package abstractions

import (
	"context" 
)

type IPromptTemplate interface {
	Render(ctx context.Context, kernel Kernel, arguments *KernelArguments) (*string, error)
}
