package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/platform"
	"github.com/opentracing/opentracing-go/log"
)

type CreatePlatformCommand struct {
	Name       string              `json:"name" validate:"required,min=3,max=50"`
	Url        string              `json:"url" validate:"required,min=3,max=50"`
	Rest       string              `json:"rest" validate:"required,min=3,max=50"`
	Tags       []string            `json:"tags"`
	Properties map[string]Property `json:"properties"`
	Secrets    map[string]Property `json:"secrets"`
}

type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreatePlatformHandler struct {
}

func (b CreatePlatformHandler) HandlerName() string {
	return "CreatePlatformHandler"
}

func (b CreatePlatformHandler) NewCommand() interface{} {
	return &CreatePlatformCommand{}
}

func (b CreatePlatformHandler) Handle(ctx context.Context, c interface{}) error {
	// if return error, it will loop and not ack
	// TODO: how to handle error, and get response?
	aux := c.(*CreatePlatformCommand)

	repository, commonHandler, err := createCommonInfra(ctx)
	if err != nil {
		log.Error(err)
		return nil
	}

	res, err := repository.GetPlatformByName(ctx, aux.Name)
	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		log.Error(err)
		return nil
	}

	if res != nil && res.Name == aux.Name {
		log.Error(fmt.Errorf("name: %s is existed", aux.Name))
		return nil
	}

	property := make(map[string]platform.Property)
	for _, v := range aux.Properties {
		property[v.Key] = platform.Property(v)
		// {
		// 	Key:      v.Key,
		// 	Value:    v.Value,
		// 	NeedMask: v.NeedMask,
		// }
	}
	err = commonHandler.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(aux.Name, aux.Url, nil,
			platform.WithPlatformProperties(property),
			platform.WithPlatformTags(aux.Tags),
		)
		return repository.Insert(ctx, *res)
	})
	if err != nil {
		log.Error(err)
		return nil
	}

	return nil
}

func createCommonInfra(ctx context.Context) (platform.IPlatformRepository, *CommonHandler, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, nil, err
	}
	repository := infra.NewPlatformRepository(client, config)

	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, nil, err
	}
	return repository, &CommonHandler{
		unitOfWork: unitOfWork,
	}, nil
}
