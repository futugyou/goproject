package mongo2struct

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Template struct {
	entityTplString        string
	repoInterfaceTplString string
	baseRepoImplTplString  string
	repoImplTplString      string
	Core                   []CoreTemplate
}

type CoreTemplate struct {
	Key string
	Tpl string
	Obj interface{}
}

func NewTemplate() *Template {
	return &Template{
		entityTplString:        entityTplString,
		repoInterfaceTplString: repoInterfaceTplString,
		baseRepoImplTplString:  base_mongorepo_TplString,
		repoImplTplString:      repoMongoImplTplString,
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

const templateName string = "mongo_struct_template"

func (t *Template) GenerateEntity(obj EntityStruct) error {
	return t.generate(t.entityTplString, fmt.Sprintf("./%s", obj.EntityFolder), fmt.Sprintf("./%s/%s.go", obj.EntityFolder, obj.FileName), obj)
}

func (t *Template) GenerateCore() error {
	for _, v := range t.Core {
		err := t.generate(v.Tpl, "./core", fmt.Sprintf("./core/%s.go", v.Key), 1)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (t *Template) GenerateBaseRepoImpl(obj interface{}) error {
	return t.generate(t.baseRepoImplTplString, "./mongorepo", "./mongorepo/respository.go", obj)
}

func (t *Template) GenerateRepository(obj RepositoryStruct) error {
	err := t.generate(t.repoImplTplString, "./mongorepo", fmt.Sprintf("./mongorepo/%s.go", obj.FileName), obj)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return t.generate(t.repoInterfaceTplString, "./repository", fmt.Sprintf("./repository/%s.go", obj.FileName), obj)
}

func (t *Template) generate(templateString string, folder string, fileName string, obj interface{}) error {
	tpl := template.Must(template.New(templateName).Funcs(template.FuncMap{
		"ToImportsList": t.toImportsList,
	}).Parse(templateString))
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		os.MkdirAll(folder, 0700)
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

func (t *Template) toImportsList(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return "import ( \n" + strings.Join(list, "\n") + "\n )"
}