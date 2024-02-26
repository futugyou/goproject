package mongo2struct

const core_repository_TplString = `
package core

import (
	"context"
	"log"
)

type IRepository[E IEntity, K any] interface {
	Insert(ctx context.Context, obj E) error
	Delete(ctx context.Context, filter []DataFilterItem) error
	GetOne(ctx context.Context, filter []DataFilterItem) (*E, error)
	Update(ctx context.Context, obj E, filter []DataFilterItem) error
	Paging(ctx context.Context, page Paging, filter []DataFilterItem) ([]E, error)
}

type DataFilter[E IEntity] func(e E) []DataFilterItem
type DataFilterItem struct {
	Key   string
	Value interface{}
}

type InsertManyResult struct {
	TabelName     string
	InsertedCount int64
	MatchedCount  int64
	ModifiedCount int64
	DeletedCount  int64
	UpsertedCount int64
}

func (i InsertManyResult) String() {
	log.Printf("table %s matched count %d \n", i.TabelName, i.MatchedCount)
	log.Printf("table %s inserted count %d \n", i.TabelName, i.InsertedCount)
	log.Printf("table %s modified count %d \n", i.TabelName, i.ModifiedCount)
	log.Printf("table %s deleted count %d \n", i.TabelName, i.DeletedCount)
	log.Printf("table %s upserted count %d \n", i.TabelName, i.UpsertedCount)
}
`
