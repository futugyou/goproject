package mongo2struct

import (
	"fmt"
	"github/go-project/tour/internal/common"

	"go.mongodb.org/mongo-driver/bson"
)

type EntityStructBuilder struct {
	EntityFolder   string
	CollectionName string
	Elements       []bson.RawElement
}

func NewEntityStructBuilder(folder string, collectionName string, elements []bson.RawElement) *EntityStructBuilder {
	return &EntityStructBuilder{
		EntityFolder:   folder,
		CollectionName: collectionName,
		Elements:       elements,
	}
}

func (b *EntityStructBuilder) Build() *common.EntityStruct {
	return &common.EntityStruct{
		EntityFolder: b.EntityFolder,
		FileName:     b.CollectionName,
		PackageName:  b.EntityFolder,
		Imports:      b.buildImports(),
		StructName:   UnderscoreToUpperCamelCase(b.CollectionName),
		Items:        b.buildItems(),
	}
}

func (b *EntityStructBuilder) buildItems() []common.EntityStructItem {
	items := make([]common.EntityStructItem, 0)
	for _, v := range b.Elements {
		itemType := convertBsontypeTogotype(v.Value())

		items = append(items, common.EntityStructItem{
			Name: UnderscoreToUpperCamelCase(v.Key()),
			Type: itemType,
			Tag:  fmt.Sprintf("`bson:\"%s\"`", v.Key()),
		})
	}
	return items
}

func (b *EntityStructBuilder) buildImports() []string {
	items := make([]string, 0)
	for _, v := range b.Elements {
		itemType := convertBsontypeTogotype(v.Value())
		items = createImports(items, itemType)
	}
	return items
}
