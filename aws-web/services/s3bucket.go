package services

import (
	"context"
	"os"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	"golang.org/x/exp/slices"
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

func (s *S3bucketService) GetS3Buckets(paging core.Paging, filter model.S3BucketFilter) []model.S3BucketViewModel {
	ctx := context.Background()
	result := make([]model.S3BucketViewModel, 0)

	f := entity.S3bucketSearchFilter{
		BucketName: filter.BucketName,
	}

	entities, err := s.repository.FilterPaging(ctx, paging, f)
	if err != nil {
		return result
	}

	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts()

	for _, entity := range entities {
		idx := slices.IndexFunc(accounts, func(c model.UserAccount) bool { return c.Id == entity.AccountId })
		bucket := model.S3BucketViewModel{
			Id:           entity.Id,
			AccountId:    entity.AccountId,
			AccountName:  accounts[idx].Alias,
			Name:         entity.Name,
			Region:       entity.Region,
			IsPublic:     entity.IsPublic,
			Policy:       entity.Policy,
			Permissions:  entity.Permissions,
			CreationDate: entity.CreationDate,
		}
		result = append(result, bucket)
	}

	return result
}
