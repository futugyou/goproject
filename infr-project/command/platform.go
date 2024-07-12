package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/platform"
)

type CreatePlatformCommand struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Rest     string            `json:"rest" validate:"required,min=3,max=50"`
	Tags     []string          `json:"tags"`
	Property map[string]string `json:"property"`
}

type CreatePlatformHandler struct {
	repository    platform.IPlatformRepository
	commonHandler CommonHandler
}

func NewCreatePlatformHandler() CreatePlatformHandler {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	repo := infra.NewPlatformRepository(client, config)
	unitOfWork, _ := infra.NewMongoUnitOfWork(client)
	return CreatePlatformHandler{
		repository: repo,
		commonHandler: CommonHandler{
			unitOfWork: unitOfWork,
		},
	}
}

func (b CreatePlatformHandler) HandlerName() string {
	return "CreatePlatformHandler"
}

func (b CreatePlatformHandler) NewCommand() interface{} {
	return &CreatePlatformCommand{}
}

func (b CreatePlatformHandler) Handle(ctx context.Context, c interface{}) error {
	aux := c.(*CreatePlatformCommand)

	var res *platform.Platform
	res, err := b.repository.GetPlatformByName(ctx, aux.Name)
	if err != nil && !strings.HasPrefix(err.Error(), "data not found") {
		return err
	}

	if res != nil && res.Name == aux.Name {
		return fmt.Errorf("name: %s is existed", aux.Name)
	}

	err = b.commonHandler.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(aux.Name, aux.Url, aux.Rest, aux.Property, aux.Tags)
		return b.repository.Insert(ctx, *res)
	})
	if err != nil {
		return err
	}

	return nil
}
