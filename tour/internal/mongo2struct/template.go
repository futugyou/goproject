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
}

func NewTemplate() *Template {
	return &Template{
		entityTplString:     entityTplString,
		repositoryTplString: "",
	}
}

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

func (t *Template) GenerateRepository(obj interface{}) error {
	return nil
}

func (t *Template) toImportsList(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return "import ( \n" + strings.Join(list, "\n") + "\n )"
}
