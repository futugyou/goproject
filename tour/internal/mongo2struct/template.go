package mongo2struct

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const entityTplString = `
package {{ .PackageName }}
{{ $import := len .Imports }}{{ if gt $import 0 }}{{ .Imports | ToImportsList }}{{ end }}
type {{ .StructName}}Entity struct {
	{{range .Items}}{{ .Name }}  {{ .Type }}  {{ .Tag }}
	{{ end }}
} 
	
func ({{ .StructName}}Entity) GetType() string {
	return "{{ .FileName }}"
}                
		 `

type Template struct {
	entityTplString     string
	repositoryTplString string
	Core                []CoreTemplate
}

type CoreTemplate struct {
	Key string
	Tpl string
}

func NewTemplate() *Template {
	return &Template{
		entityTplString:     entityTplString,
		repositoryTplString: "",
		Core: []CoreTemplate{{
			Key: "entity",
			Tpl: core_entity_TplString,
		}, {
			Key: "repository",
			Tpl: core_repository_TplString,
		}, {
			Key: "page",
			Tpl: core_page_TplString,
		}},
	}
}

const core_entity_TplString = `
package core

type IEntity interface {
	GetType() string
}
`

const core_page_TplString = `
package core

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

func (t *Template) GenerateEntity(obj EntityStruct) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToImportsList": t.toImportsList,
	}).Parse(t.entityTplString))
	if _, err := os.Stat(fmt.Sprintf("./%s/%s.go", obj.EntityFolder, obj.FileName)); os.IsNotExist(err) {
		os.MkdirAll(fmt.Sprintf("./%s", obj.EntityFolder), 0700)
	}

	f, err := os.OpenFile(fmt.Sprintf("./%s/%s.go", obj.EntityFolder, obj.FileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = tpl.Execute(f, obj)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (t *Template) GenerateCore() error {
	for _, v := range t.Core {
		tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
			"ToImportsList": t.toImportsList,
		}).Parse(v.Tpl))
		if _, err := os.Stat(fmt.Sprintf("./core/%s.go", v.Key)); os.IsNotExist(err) {
			os.MkdirAll("./core", 0700)
		}

		f, err := os.OpenFile(fmt.Sprintf("./core/%s.go", v.Key), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = tpl.Execute(f, 1)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (t *Template) GenerateRepository(obj interface{}) error {
	return nil
}

func (t *Template) toImportsList(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return "import ( \n" + strings.Join(list, "\n") + "\n )"
}
