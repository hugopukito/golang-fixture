package funcs

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

type Temp struct {
	Test string
}

var typeRegistry = make(map[string]reflect.Type)

func InitLocalStructs(pkgName string) {
	structNames, err := getAllStructsInPackage(pkgName)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(structNames)
	// structType := reflect.TypeOf("Temp")
	// structInstance := reflect.New(structType.Elem()).Interface()
	// myTypes := []interface{}{
	// 	//Dog{},
	// 	structInstance,
	// }

	// for _, v := range myTypes {
	// 	fmt.Printf("%T \n", v)
	// 	typeRegistry[fmt.Sprintf("%T", v)] = reflect.TypeOf(v)
	// }
}

func StructAssign(entityStructName string) {
	// instance := createStructIfPossible(entityStructName)
	// value := reflect.ValueOf(instance)
	// typ := reflect.TypeOf(instance)

	// for i := 0; i < typ.NumField(); i++ {
	// 	field := typ.Field(i)
	// 	fieldValue := value.Field(i).Interface()

	// 	fmt.Printf("Field Name: %s\n", field.Name)
	// 	fmt.Printf("Field Type: %s\n", field.Type)
	// 	fmt.Printf("Field Value: %v\n\n", fieldValue)
	// }
}

func CreateStructIfPossible(pkgName string, structName string) interface{} {
	v := reflect.New(typeRegistry[pkgName+"."+structName]).Elem()
	return v.Interface()
}

func getAllStructsInPackage(pkgName string) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("Error getting current working directory: " + err.Error())
	}
	structDir := wd + "/" + pkgName

	files, err := os.ReadDir(structDir)
	if err != nil {
		return nil, errors.New("Error reading directory: " + err.Error())
	}

	structNames := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(structDir, file.Name())

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					structName := typeSpec.Name.Name
					structNames = append(structNames, structName)

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					for _, field := range structType.Fields.List {
						for _, fieldName := range field.Names {
							structNames = append(structNames, structName+"."+fieldName.Name)
						}
					}
				}
			}
		}
	}
	return structNames, nil
}
