package services

import (
	"os"

	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
)

type S3bucketService struct {
	repository     repository.IS3bucketRepository
	itemRepository repository.IS3bucketItemRepository
}

func NewS3bucketService() *S3bucketService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &S3bucketService{
		repository:     mongorepo.NewS3bucketRepository(config),
		itemRepository: mongorepo.NewS3bucketItemRepository(config),
	}
}
