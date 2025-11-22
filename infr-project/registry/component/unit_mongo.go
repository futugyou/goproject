package component

import (
	"context"

	"github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/infr-project/registry/options"
	"github.com/futugyou/infr-project/registry/unit"
)

func init() {
	unit.DefaultRegistry.RegisterComponent(func(ctx context.Context, option options.Options) domain.UnitOfWork {
		mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
		if err != nil {
			return nil
		}

		unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
		return unitOfWork
	}, "mongo")
}
