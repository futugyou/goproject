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
						// get field.Type as string
						var typeNameBuf bytes.Buffer
						err := printer.Fprint(&typeNameBuf, fset, field.Type)
						if err != nil {
							log.Fatalf("failed printing %s", err)
							continue
						}

						fieldName := ""
						fieldTypeName := typeNameBuf.String()
						fieldTag := ""
						if field.Tag != nil {
							fieldTag = field.Tag.Value
						}

						// if nestedStruct, ok := field.Type.(*ast.StructType); ok {
						// 	for _, nestedField := range nestedStruct.Fields.List {
						// 		nestedFieldName := nestedField.Names[0].Name
						// 		var typeNameBuf2 bytes.Buffer
						// 		printer.Fprint(&typeNameBuf2, fset, nestedField.Type)
						// 		fmt.Printf("  1%s: %s\n", nestedFieldName, typeNameBuf2.String())
						// 	}
						// }

						// if arrayType, ok := field.Type.(*ast.ArrayType); ok {
						// 	var typeNameBuf2 bytes.Buffer
						// 	printer.Fprint(&typeNameBuf2, fset, arrayType.Elt)
						// 	fmt.Printf("  2 : %s\n", typeNameBuf2.String())

						// }

						if len(field.Names) == 0 {
							fieldName = fieldTypeName
						} else {
							fieldName = field.Names[0].Name
						}
						fls = append(fls, FieldInfo{
							Name:     fieldName,
							TypeName: fieldTypeName,
							Tag:      fieldTag,
						})
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
