package component

import (
	"context"

	"github.com/futugyou/infr-project/domain"
	memory "github.com/futugyou/infr-project/infrastructure_memory"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/unit"
)

func init() {
	unit.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) domain.IUnitOfWork {
		return memory.NewMemoryUnitOfWork()
	}, "memory")
}
