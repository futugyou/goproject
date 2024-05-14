package main

import (
	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/resource"

	"github.com/google/wire"
)

func ProvideSourcer() application.IEventSourcingService[resource.IResourceEvent, *resource.Resource] {
	return nil //application.NewEventSourcer[resource.IResourceEvent, *resource.Resource]()
}

// https://github.com/google/wire/pull/360#issuecomment-1141376353
// it is not work
// panic: unhandled AST node: *ast.IndexListExpr [recovered]
func InitializeResourceService() *application.ResourceService {
	wire.Build(application.NewResourceService, ProvideSourcer)
	return &application.ResourceService{}
}
