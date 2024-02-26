package mongo2struct

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Template struct {
	entityTplString       string
	repositoryTplString   string
	baseRepoImplTplString string
	Core                  []CoreTemplate
}

type CoreTemplate struct {
	Key string
	Tpl string
	Obj interface{}
}

func NewTemplate() *Template {
	return &Template{
		entityTplString:       entityTplString,
		repositoryTplString:   "",
		baseRepoImplTplString: base_mongorepo_TplString,
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
	tpl := template.Must(template.New(templateName).Funcs(template.FuncMap{
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
		tpl := template.Must(template.New(templateName).Funcs(template.FuncMap{
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

func (t *Template) GenerateBaseRepoImpl(obj interface{}) error {
	tpl := template.Must(template.New(templateName).Parse(t.baseRepoImplTplString))
	if _, err := os.Stat("./mongorepo/respository.go"); os.IsNotExist(err) {
		os.MkdirAll("./mongorepo", 0700)
	}

	f, err := os.OpenFile("./mongorepo/respository.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

func (t *Template) GenerateRepository(obj interface{}) error {
	return nil
}

func (t *Template) toImportsList(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return "import ( \n" + strings.Join(list, "\n") + "\n )"
}
