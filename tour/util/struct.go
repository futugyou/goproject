package util

import (
	"bytes"
	"fmt"
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
	FilePath        string
	structInfoCache []StructInfo
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

func NewASTManager(filePath string) *ASTManager {
	return &ASTManager{
		FilePath:        filePath,
		structInfoCache: []StructInfo{},
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

	if fileInfo.IsDir() {
		packages, err := parser.ParseDir(fset, m.FilePath, nil, 0)
		if err != nil {
			return nil, err
		}
		for _, pack := range packages {
			for _, file := range pack.Files {
				ast.Inspect(file, m.astInspectFunc(fset, &structs, file.Name.Name))
			}
		}
	} else {
		file, err := parser.ParseFile(fset, m.FilePath, nil, 0)
		if err != nil {
			return nil, err
		}
		ast.Inspect(file, m.astInspectFunc(fset, &structs, file.Name.Name))
	}
	m.structInfoCache = structs
	return
}

func (m *ASTManager) astInspectFunc(fset *token.FileSet, structs *[]StructInfo, packageName string) func(ast.Node) bool {
	return func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			fls := make([]FieldInfo, 0)
			v := x.Type.(*ast.StructType)
			for _, field := range v.Fields.List {
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

				fieldName := ""
				fieldTypeName := typeNameBuf.String()
				fieldTag := ""
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
			*structs = append(*structs, StructInfo{
				PackageName: packageName,
				StructName:  x.Name.Name,
				FieldInfos:  fls,
			})
		}

		return true
	}
}

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

						fieldName := ""
						fieldTypeName := typeNameBuf.String()
						fieldTag := ""
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

func stringToReflectType(t string, structs []StructInfo) (reflect.Type, error) {
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
				keyType, err := stringToReflectType(maptypeSplit[0], structs)
				if err != nil {
					return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
				}
				valueType, err := stringToReflectType(maptypeSplit[1], structs)
				if err != nil {
					return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
				}
				return reflect.MapOf(keyType, valueType), nil
			}
		}
	}

	// handle array
	if arrayType, ok := strings.CutPrefix(t, "[]"); ok {
		keyType, err := stringToReflectType(arrayType, structs)
		if err != nil {
			return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
		}
		return reflect.ArrayOf(0, keyType), nil
	}

	return GetReflectTypeFromStructInfo(t, structs)
}

func (m *ASTManager) GetReflectTypeByName(structName string) (reflect.Type, error) {
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
		ty, err := stringToReflectType(v.TypeName, structs)
		if err != nil {
			log.Println(err)
			continue
		}
		anonymous := false
		if v.Name == v.TypeName {
			anonymous = true
		}
		fields = append(fields, reflect.StructField{
			Name:      v.Name,
			Type:      ty,
			Tag:       reflect.StructTag(v.Tag),
			Anonymous: anonymous,
		})
	}

	return reflect.StructOf(fields), nil
}

func (m *ASTManager) GetAllReflectType() ([]reflect.Type, error) {
	structs, err := m.GetStructInfo()
	if err != nil {
		return nil, err
	}

	result := make([]reflect.Type, 0)
	for _, info := range structs {
		fields := make([]reflect.StructField, 0)
		for _, v := range info.FieldInfos {
			ty, err := stringToReflectType(v.TypeName, structs)
			if err != nil {
				log.Println(err)
				continue
			}
			anonymous := false
			if v.Name == v.TypeName {
				anonymous = true
			}
			fields = append(fields, reflect.StructField{
				Name:      v.Name,
				Type:      ty,
				Tag:       reflect.StructTag(v.Tag),
				Anonymous: anonymous,
			})
		}
		reflectType := reflect.StructOf(fields)

		result = append(result, reflectType)
	}

	return result, nil
}

func GetReflectTypeFromStructInfo(t string, structs []StructInfo) (reflect.Type, error) {
	var info *StructInfo
	for _, s := range structs {
		if t == s.StructName {
			info = &StructInfo{
				PackageName: s.PackageName,
				StructName:  s.StructName,
				FieldInfos:  s.FieldInfos,
			}
			break
		}
	}

	if info == nil || len(info.FieldInfos) == 0 {
		return nil, fmt.Errorf("%s can not convert to reflect.Type", t)
	}

	fields := make([]reflect.StructField, 0)
	for _, v := range info.FieldInfos {
		ty, err := stringToReflectType(v.TypeName, structs)
		if err != nil {
			log.Println(err)
			continue
		}
		anonymous := false
		if v.Name == v.TypeName {
			anonymous = true
		}
		fields = append(fields, reflect.StructField{
			Name:      v.Name,
			Type:      ty,
			Tag:       reflect.StructTag(v.Tag),
			Anonymous: anonymous,
		})
	}

	return reflect.StructOf(fields), nil
}

func CreateInstanceByStructInfos(structs []StructInfo) []reflect.Type {
	result := make([]reflect.Type, 0)
	for _, info := range structs {
		fields := make([]reflect.StructField, 0)
		for _, v := range info.FieldInfos {
			ty, err := stringToReflectType(v.TypeName, structs)
			if err != nil {
				log.Println(err)
				continue
			}
			anonymous := false
			if v.Name == v.TypeName {
				anonymous = true
			}
			fields = append(fields, reflect.StructField{
				Name:      v.Name,
				Type:      ty,
				Tag:       reflect.StructTag(v.Tag),
				Anonymous: anonymous,
			})
		}
		reflectType := reflect.StructOf(fields)

		result = append(result, reflectType)
	}

	return result
}
