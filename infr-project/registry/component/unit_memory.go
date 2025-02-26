package component

import (
	"github.com/futugyou/infr-project/domain"
	memory "github.com/futugyou/infr-project/infrastructure_memory"
	"github.com/futugyou/infr-project/options"
	"github.com/futugyou/infr-project/registry/unit"
)

func init() {
	unit.DefaultRegistry.RegisterComponent(func(option options.Options) domain.IUnitOfWork {
		return memory.NewMemoryUnitOfWork()
	}, func(option options.Options) domain.IUnitOfWorkAsync {
		return memory.NewMemoryUnitOfWork()
	}, "memory")
}
