package mongo2struct

import (
	"fmt"
	"os"
	"text/template"
)

const structTpl = `
package {{ .PackageName }}

type {{ .StructName}}Entity struct {
	{{range .Items}}{{ .Name }}  {{ .Type }}  {{ .Tag }}
	{{ end }}
} 
	
func ({{ .StructName}}Entity) GetType() string {
	return "{{ .PackageName }}"
}                
		 `

type StructTemplate struct {
	structTpl string
	obj       interface{}
}

func NewStructTemplate(obj Struct) *StructTemplate {
	return &StructTemplate{
		structTpl: structTpl,
		obj:       obj,
	}
}

func (t *StructTemplate) Generate() error {
	tpl := template.Must(template.New("sql2struct").Parse(t.structTpl))

	err := tpl.Execute(os.Stdout, t.obj)
	if err != nil {
		return err
	}
	fmt.Println()
	return nil
}
