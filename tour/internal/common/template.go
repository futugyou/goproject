package common

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

func NewDefaultTemplate(baseRepoTplString string) *Template {
	return &Template{
		entityTplString:        entityTplString,
		repoInterfaceTplString: repoInterfaceTplString,
		baseRepoImplTplString:  baseRepoTplString,
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

func (t *Template) GenerateCore(obj CoreConfig) error {
	for _, v := range t.Core {
		err := t.generate(v.Tpl, fmt.Sprintf("./%s", obj.Folder), fmt.Sprintf("./%s/%s.go", obj.Folder, v.Key), obj)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (t *Template) GenerateEntity(obj EntityStruct) error {
	return t.generate(t.entityTplString, fmt.Sprintf("./%s", obj.EntityFolder), fmt.Sprintf("./%s/%s.go", obj.EntityFolder, obj.FileName), obj)
}

func (t *Template) GenerateBaseRepoImpl(obj BaseRepoImplConfig) error {
	return t.generate(t.baseRepoImplTplString, fmt.Sprintf("./%s", obj.Folder), fmt.Sprintf("./%s/%s.go", obj.Folder, obj.FileName), obj.TemplateObj)
}

func (t *Template) GenerateRepository(obj RepositoryStruct) error {
	err := t.generate(t.repoImplTplString, fmt.Sprintf("./%s", obj.Folder), fmt.Sprintf("./%s/%s.go", obj.Folder, obj.FileName), obj)
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
