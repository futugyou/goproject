package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IS3bucketRepository interface {
	core.IRepository[entity.S3bucketEntity, string]
	DeleteAll(ctx context.Context) error
	FilterPaging(ctx context.Context, page core.Paging, filter entity.S3bucketSearchFilter) ([]*entity.S3bucketEntity, error)
}

type IS3bucketItemRepository interface {
	core.IRepository[entity.S3bucketItemEntity, string]
	DeleteByBucketName(ctx context.Context, bucketName string) error
	FilterPaging(ctx context.Context, page core.Paging, filter entity.S3bucketSearchFilter) ([]*entity.S3bucketItemEntity, error)
}
