package util

import (
	"bytes"
	"fmt"
	"github/go-project/tour/internal/word"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type ASTManager struct {
	FilePath         string
	structInfoCache  []StructInfo
	reflectTypeCache map[string]reflect.Type
}

type StructInfo struct {
	PackageName string
	StructName  string
	Doc         string
	Comment     string
	FieldInfos  []FieldInfo
}

type FieldInfo struct {
	Name     string
	TypeName string
	Tag      string
	Doc      string
	Comment  string
}

func (f *FieldInfo) String() string {
	return fmt.Sprintf("%s %s %s", f.Name, f.TypeName, f.Tag)
}

func NewASTManager(filePath string) *ASTManager {
	return &ASTManager{
		FilePath:         filePath,
		structInfoCache:  []StructInfo{},
		reflectTypeCache: map[string]reflect.Type{},
	}
}

func (m *ASTManager) GetStructInfo() (structs []StructInfo, err error) {
	if len(m.structInfoCache) > 0 {
		return m.structInfoCache, nil
	}

	fset := token.NewFileSet()
	fileInfo, err := os.Stat(m.FilePath)
	if err != nil {
		return nil, err
	}
	var currentGenDecl *ast.GenDecl // Keep track of the current GenDecl
	if fileInfo.IsDir() {
		packages, err := parser.ParseDir(fset, m.FilePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, pack := range packages {
			for _, file := range pack.Files {
				ast.Inspect(file, m.astInspectFunc(fset, &structs, file.Name.Name, currentGenDecl))
			}
		}
	} else {
		file, err := parser.ParseFile(fset, m.FilePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		ast.Inspect(file, m.astInspectFunc(fset, &structs, file.Name.Name, currentGenDecl))
	}
	m.structInfoCache = structs
	return
}

func (m *ASTManager) astInspectFunc(fset *token.FileSet, structs *[]StructInfo, packageName string, currentGenDecl *ast.GenDecl) func(ast.Node) bool {
	return func(n ast.Node) bool {
		// switch x := n.(type) {
		// case *ast.FuncDecl:
		// 	fmt.Printf("%s:\tFuncDecl %s\t%s\n", fset.Position(n.Pos()), x.Name, x.Doc.Text())
		// case *ast.TypeSpec:
		// 	fmt.Printf("%s:\tTypeSpec %s\t%s\n", fset.Position(n.Pos()), x.Name, x.Doc.Text())
		// case *ast.Field:
		// 	fmt.Printf("%s:\tField %s\t%s\n", fset.Position(n.Pos()), x.Names, x.Doc.Text())
		// case *ast.GenDecl:
		// 	fmt.Printf("%s:\tGenDecl %s\n", fset.Position(n.Pos()), x.Doc.Text())
		// }

		switch x := n.(type) {
		case *ast.GenDecl:
			currentGenDecl = x
		case *ast.TypeSpec:
			structDoc := strings.TrimSpace(x.Doc.Text())
			structComment := strings.TrimSpace(x.Comment.Text())

			if x.Doc == nil {
				structDoc = strings.TrimSpace(currentGenDecl.Doc.Text())
			}

			fls := make([]FieldInfo, 0)
			structType := x.Type.(*ast.StructType)

			for _, field := range structType.Fields.List {
				// if mapType, ok := field.Type.(*ast.MapType); ok {
				// 	fmt.Println(mapType.Key, mapType.Value)
				// }

				// get field.Type as string
				var typeNameBuf bytes.Buffer
				err := printer.Fprint(&typeNameBuf, fset, field.Type)
				if err != nil {
					log.Fatalf("failed printing %s", err)
					continue
				}

				fieldTypeName := typeNameBuf.String()
				fieldName := ""
				fieldTag := ""
				fieldDoc := strings.TrimSpace(field.Doc.Text())
				fieldComment := strings.TrimSpace(field.Comment.Text())
				if len(field.Names) == 0 {
					fieldName = fieldTypeName
				} else {
					fieldName = field.Names[0].Name
				}

				if field.Tag != nil {
					fieldTag = strings.ReplaceAll(field.Tag.Value, "`", "")
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

				fls = append(fls, FieldInfo{
					Name:     fieldName,
					TypeName: fieldTypeName,
					Tag:      fieldTag,
					Doc:      fieldDoc,
					Comment:  fieldComment,
				})
			}
			*structs = append(*structs, StructInfo{
				PackageName: packageName,
				StructName:  x.Name.Name,
				Doc:         structDoc,
				Comment:     structComment,
				FieldInfos:  fls,
			})
		}

		return true
	}
}

func (m *ASTManager) stringToReflectType(t string) (reflect.Type, error) {
	switch t {
	case "string":
		return reflect.TypeOf(""), nil
	case "bool":
		return reflect.TypeOf(false), nil
	case "int", "int64", "int16", "int8", "int32", "uint", "uint16", "uint32", "uint64", "uint8":
		return reflect.TypeOf(int64(0)), nil
	case "float32", "float64":
		return reflect.TypeOf(float64(0)), nil
	case "time.Time":
		return reflect.TypeOf(time.Time{}), nil
	}

	// handle map, only handle map[xxx]xxx, not map[xxx]map[xxx]xxx...
	if strings.Contains(t, "map[") {
		if maptype, ok := strings.CutPrefix(t, "map["); ok {
			if maptypeSplit := strings.Split(maptype, "]"); len(maptypeSplit) == 2 {
				keyType, err := m.stringToReflectType(maptypeSplit[0])
				if err != nil {
					return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
				}
				valueType, err := m.stringToReflectType(maptypeSplit[1])
				if err != nil {
					return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
				}
				return reflect.MapOf(keyType, valueType), nil
			}
		}
	}

	// handle array
	if arrayType, ok := strings.CutPrefix(t, "[]"); ok {
		keyType, err := m.stringToReflectType(arrayType)
		if err != nil {
			return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
		}
		return reflect.ArrayOf(0, keyType), nil
	}

	return m.GetReflectTypeByName(t)
}

func (m *ASTManager) GetReflectTypeByName(structName string) (reflect.Type, error) {
	if rt, ok := m.reflectTypeCache[structName]; ok {
		return rt, nil
	}

	structs, err := m.GetStructInfo()
	if err != nil {
		return nil, err
	}

	var info *StructInfo
	for _, s := range structs {
		if structName == s.StructName {
			info = &StructInfo{
				PackageName: s.PackageName,
				StructName:  s.StructName,
				FieldInfos:  s.FieldInfos,
			}
			break
		}
	}

	if info == nil || len(info.FieldInfos) == 0 {
		return nil, fmt.Errorf("%s can not convert to reflect.Type", structName)
	}

	fields := make([]reflect.StructField, 0)
	for _, v := range info.FieldInfos {
		var ty reflect.Type
		var ok bool
		if ty, ok = m.reflectTypeCache[v.TypeName]; !ok {
			ty, err = m.stringToReflectType(v.TypeName)
			if err != nil {
				log.Println(err)
				continue
			}
			m.reflectTypeCache[v.TypeName] = ty
		}

		anonymous := false
		if v.Name == v.TypeName {
			anonymous = true
		}

		tag := reflect.StructTag(v.Tag)
		if _, ok := tag.Lookup("description"); !ok {
			if len(v.Comment) > 0 {
				tag = reflect.StructTag(fmt.Sprintf("%s description:\"%s\"", v.Tag, v.Comment))
			} else if len(v.Doc) > 0 {
				tag = reflect.StructTag(fmt.Sprintf("%s description:\"%s\"", v.Tag, v.Doc))
			}
		}

		if text, ok := tag.Lookup("bson"); ok {
			tag = reflect.StructTag(fmt.Sprintf("%s json:\"%s\"", v.Tag, text))
		}

		if _, ok := tag.Lookup("json"); !ok {
			tag = reflect.StructTag(fmt.Sprintf("%s json:\"%s\"", v.Tag, word.CamelCaseToUnderscore(v.Name)))
		}

		fields = append(fields, reflect.StructField{
			Name:      v.Name,
			Type:      ty,
			Tag:       tag,
			Anonymous: anonymous,
		})
	}

	rt := reflect.StructOf(fields)
	m.reflectTypeCache[structName] = rt

	return rt, nil
}

func (m *ASTManager) GetAllReflectType() (map[string]reflect.Type, error) {
	if len(m.reflectTypeCache) > 0 {
		return m.reflectTypeCache, nil
	}
	result := map[string]reflect.Type{}
	structs, err := m.GetStructInfo()
	if err != nil {
		return nil, err
	}

	for _, info := range structs {
		reflectType, err := m.GetReflectTypeByName(info.StructName)
		if err != nil {
			return nil, err
		}
		result[info.StructName] = reflectType
	}

	return result, nil
}
