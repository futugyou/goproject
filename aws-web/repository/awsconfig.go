package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IAwsConfigRepository interface {
	core.IRepository[entity.AwsConfigEntity, string]
	BulkWrite(ctx context.Context, entities []entity.AwsConfigEntity) error
}

type IAwsConfigRelationshipRepository interface {
	core.IRepository[entity.AwsConfigRelationshipEntity, string]
	BulkWrite(ctx context.Context, entities []entity.AwsConfigRelationshipEntity) error
}
