package services

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
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

func (s *S3bucketService) GetS3Buckets(ctx context.Context, paging core.Paging, filter model.S3BucketFilter) []model.S3BucketViewModel {
	result := make([]model.S3BucketViewModel, 0)

	f := entity.S3bucketSearchFilter{
		BucketName: filter.BucketName,
	}

	entities, err := s.repository.FilterPaging(ctx, paging, f)
	if err != nil {
		return result
	}

	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts(ctx)

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

func (s *S3bucketService) GetS3BucketItems(ctx context.Context, filter model.S3BucketItemFilter) []model.S3BucketItemViewModel {
	result := make([]model.S3BucketItemViewModel, 0)
	accountService := NewAccountService()
	account := accountService.GetAccountByID(ctx, filter.AccountId)
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	output, err := s.ListItems(ctx, filter.BucketName, filter.Perfix, filter.Del)
	if err != nil {
		return nil
	}

	for _, obj := range output.Contents {
		i := model.S3BucketItemViewModel{
			Id:           *obj.Key,
			BucketName:   filter.BucketName,
			Key:          *obj.Key,
			Size:         *obj.Size,
			CreationDate: *obj.LastModified,
			IsDirectory:  false,
		}
		result = append(result, i)
	}

	for _, obj := range output.CommonPrefixes {
		i := model.S3BucketItemViewModel{
			Id:           *obj.Prefix,
			BucketName:   filter.BucketName,
			Key:          *obj.Prefix,
			Size:         0,
			CreationDate: time.Time{},
			IsDirectory:  true,
		}
		result = append(result, i)
	}
	return result
}

func (s *S3bucketService) GetS3BucketFile(ctx context.Context, filter model.S3BucketFileFilter) (*s3.GetObjectOutput, error) {
	accountService := NewAccountService()
	account := accountService.GetAccountByID(ctx, filter.AccountId)
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)

	return s.GetS3Object(ctx, filter.BucketName, filter.FileName)
}

func (s *S3bucketService) GetS3FileUrl(ctx context.Context, filter model.S3BucketFileFilter) model.S3BucketUrlViewModel {
	accountService := NewAccountService()
	account := accountService.GetAccountByID(ctx, filter.AccountId)
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	url := s.PresignGetObject(ctx, filter.BucketName, filter.FileName)

	return model.S3BucketUrlViewModel{Url: url}
}
