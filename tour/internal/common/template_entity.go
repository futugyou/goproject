package common

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
