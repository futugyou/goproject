package util

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
)

func GetStructsFromFolder(filePath string) (structs []StructInfo, err error) {
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, filePath, nil, 0)
	for _, pack := range packages {
		for _, file := range pack.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					fls := make([]FieldInfo, 0)
					v := x.Type.(*ast.StructType)
					for _, field := range v.Fields.List {
						for _, name := range field.Names {
							// get field.Type as string
							var typeNameBuf bytes.Buffer
							err := printer.Fprint(&typeNameBuf, fset, field.Type)
							if err != nil {
								log.Fatalf("failed printing %s", err)
								continue
							}
							fls = append(fls, FieldInfo{
								Name:     name.Name,
								TypeName: typeNameBuf.String(),
								Tag:      field.Tag.Value,
							})
						}

					}
					structs = append(structs, StructInfo{
						PackageName: pack.Name,
						StructName:  x.Name.Name,
						FieldInfos:  fls,
					})
				}

				return true
			})
		}
	}

	return
}

type StructInfo struct {
	PackageName string
	StructName  string
	FieldInfos  []FieldInfo
}

type FieldInfo struct {
	Name     string
	TypeName string
	Tag      string
}

func (f *FieldInfo) String() string {
	return fmt.Sprintf("%s %s %s", f.Name, f.TypeName, f.Tag)
}
