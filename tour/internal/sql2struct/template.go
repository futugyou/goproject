package sql2struct

import (
	"fmt"
	"github/go-project/tour/internal/word"
	"os"
	"text/template"
)

const structTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}  //{{ $length :=len .Comment}}{{ if gt $length 0 }} {{ .Comment }}{{else}}{{ .Name | ToCamelCase }}{{ end }}
  {{ $typeLen := len .Type }}{{ if gt $typeLen 0 }}{{ .Name | ToCamelCase }} {{ .Type }} {{ .Tag  | ToUnderlineCase}}{{ else }}{{ .Name | ToCamelCase }}{{ end }}
{{ end }}} 

func (model {{ .TableName | ToCamelCase }}) TableName() string{
  return "{{ .TableName }}"
}                 
 	`

type StructTemplate struct {
	structTpl string
}

type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type structTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: structTpl}
}

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tpCols := make([]*StructColumn, 0, len(tbColumns))
	for _, col := range tbColumns {
		tpCols = append(tpCols, &StructColumn{
			Name:    col.ColumnName,
			Type:    col.DataType,
			Tag:     fmt.Sprintf("`json:"+"%s"+"`", col.ColumnName),
			Comment: col.ColumnComment,
		})
	}
	return tpCols
}

func (t *StructTemplate) Generate(tablename string, tplcol []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase":     word.UnderscoreToUpperCamelCase,
		"ToUnderlineCase": word.CamelCaseToUnderscore,
	}).Parse(t.structTpl))
	tplDB := structTemplateDB{
		TableName: tablename,
		Columns:   tplcol,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}
	return nil
}
