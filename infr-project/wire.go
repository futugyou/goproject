package main

import (
	"github.com/futugyou/infr-project/application"
	"github.com/google/wire"
)

func InitializeResourceService(phrase string) (*application.ResourceService, error) {
	wire.Build(application.NewResourceService())
	return nil, nil
}
