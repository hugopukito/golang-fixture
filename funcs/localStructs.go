package funcs

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

var structMap = make(map[string]map[string]string)

func InitLocalStructs(pkgName string) {
	err := getAllStructsInPackage(pkgName)
	if err != nil {
		log.Panicln(err)
	}
}

func GetFieldsFromStructName(structName string) (map[string]string, bool) {
	value, ok := structMap[structName]
	return value, ok
}

func getAllStructsInPackage(pkgName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.New("Error getting current working directory: " + err.Error())
	}
	structDir := wd + "/" + pkgName

	files, err := os.ReadDir(structDir)
	if err != nil {
		return errors.New("Error reading directory: " + err.Error())
	}

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
					fieldMap := make(map[string]string)

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					for _, field := range structType.Fields.List {
						fieldName := ""
						fieldType := ""

						for _, fieldNameIdent := range field.Names {
							fieldName = fieldNameIdent.Name
						}

						switch fieldTypeExpr := field.Type.(type) {
						case *ast.Ident:
							fieldType = fieldTypeExpr.Name
						case *ast.SelectorExpr:
							fieldType = fieldTypeExpr.Sel.Name
						}

						fieldMap[fieldName] = fieldType
					}

					structMap[structName] = fieldMap
				}
			}
		}
	}
	return nil
}
