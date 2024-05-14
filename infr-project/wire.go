package main

import (
	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/services"
	"github.com/google/wire"
)

func ProvideSourcer() application.IEventSourcingService[services.IResourceEvent, *services.Resource] {
	return nil //application.NewEventSourcer[services.IResourceEvent, *services.Resource]()
}

// https://github.com/google/wire/pull/360#issuecomment-1141376353
// it is not work
// panic: unhandled AST node: *ast.IndexListExpr [recovered]
func InitializeResourceService() *application.ResourceService {
	wire.Build(application.NewResourceService, ProvideSourcer)
	return &application.ResourceService{}
}
