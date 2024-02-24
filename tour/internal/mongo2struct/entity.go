package mongo2struct

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type EntityStruct struct {
	EntityFolder string
	FileName     string
	PackageName  string
	Imports      []string
	StructName   string
	Items        []EntityStructItem
}

type EntityStructItem struct {
	Name string
	Type string
	Tag  string
}

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

func (b *EntityStructBuilder) Build() *EntityStruct {
	return &EntityStruct{
		EntityFolder: b.EntityFolder,
		FileName:     b.CollectionName,
		PackageName:  b.EntityFolder,
		Imports:      b.buildImports(),
		StructName:   UnderscoreToUpperCamelCase(b.CollectionName),
		Items:        b.buildItems(),
	}
}

func (b *EntityStructBuilder) buildItems() []EntityStructItem {
	items := make([]EntityStructItem, 0)
	for _, v := range b.Elements {
		itemType := convertBsontypeTogotype(v.Value())

		items = append(items, EntityStructItem{
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
