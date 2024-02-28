package mongo2struct

type CoreConfig struct {
	PackageName string
	Folder      string
}

const core_entity_TplString = `
package {{ .PackageName }}

type IEntity interface {
	GetType() string
}
`

const core_page_TplString = `
package {{ .PackageName }}

type Paging struct {
	Page      int64
	Limit     int64
	SortField string
	Direct    SortDirect
}

const ASC sortDirect = "ASC"
const DESC sortDirect = "DESC"

type SortDirect interface {
	privateSortDirect()
	String() string
}

type sortDirect string

func (c sortDirect) privateSortDirect() {}
func (c sortDirect) String() string {
	return string(c)
}
`

const core_repository_TplString = `
package {{ .PackageName }}

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
