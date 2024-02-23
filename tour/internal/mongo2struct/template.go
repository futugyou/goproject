package mongo2struct

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const structTpl = `
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

type StructTemplate struct {
	structTpl string
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{
		structTpl: structTpl,
	}
}

func (t *StructTemplate) Generate(obj Struct) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToImportsList": ToImportsList,
	}).Parse(t.structTpl))
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

func ToImportsList(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return "import ( \n" + strings.Join(list, "\n") + "\n )"
}
