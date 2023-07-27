package funcs

import (
	"errors"
	"fixture/color"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var structMap = make(map[string]map[string]string)
var specialTypes map[string]func(any) bool
var randomRegex *regexp.Regexp

func init() {
	specialTypes = make(map[string]func(any) bool)
	addFuncsToSpecialTypes()
	compileRegex()
}

func compileRegex() {
	var err error
	randomRegex, err = regexp.Compile(`\{random\{([^}]*)\}\}`)
	if err != nil {
		log.Println("Failed to compile regular expression:", err)
	}
}

func addFuncsToSpecialTypes() {
	isUUID := func(obj any) bool {
		switch val := obj.(type) {
		case string:
			_, err := uuid.Parse(val)
			return err == nil
		default:
			return false
		}
	}
	specialTypes["uuid.UUID"] = isUUID

	isTime := func(obj any) bool {
		switch val := obj.(type) {
		case time.Time:
			return true
		case string:
			_, err := time.Parse("2006-01-02 15:04:05", val)
			return err == nil
		default:
			return false
		}
	}
	specialTypes["time.Time"] = isTime
}

func InitLocalStructs(pkgName string) {
	err := getAllStructsInPackage(pkgName)
	if err != nil {
		log.Panicln(color.Red + "getAllStructsInPackage err: " + err.Error() + color.Reset)
	}
}

func GetLocalStructByName(structName string) (map[string]string, bool) {
	localStruct, ok := structMap[structName]
	return localStruct, ok
}

func CheckEntityOfStructIsValid(structName string, entity map[string]any, entityName string) bool {
	fieldsAndTypes, ok := structMap[structName]
	if !ok {
		fmt.Println(color.Red+"Unknown struct ->", color.Orange, structName+"..."+color.Reset)
		return false
	}
	for field, value := range entity {
		localType, ok := fieldsAndTypes[field]
		if !ok {
			fmt.Println(color.Red+"Type "+color.Orange+field+color.Red+" doesn't exit for entity ->", color.Orange, entityName+"..."+color.Reset)
			return false
		}
		fixtureType := reflect.TypeOf(value).Name()
		if localType != fixtureType {
			if _func, ok := specialTypes[localType]; !(ok && _func(value)) {
				if _, isString := value.(string); !isString {
					fmt.Println(color.Red+"local type: "+color.Orange+localType+color.Red+" doesn't match with entity type: "+color.Orange+fixtureType+color.Red+" on field: "+color.Orange+field+color.Red+" and unknown type value for entity ->", entityName+color.Reset)
				} else {
					if randomRegex.MatchString(value.(string)) {
						return true
					} else {
						fmt.Println(color.Red+"local type: "+color.Orange+localType+color.Red+" doesn't match with entity type: "+color.Orange+fixtureType+color.Red+" on field and value: "+color.Orange+field+": "+value.(string)+color.Red+" for entity ->", entityName+color.Reset)
					}
				}
				return false
			}
		}
	}

	return true
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
				log.Panicln(color.Red + "getAllStructsInPackage ParseFile err: " + err.Error() + color.Reset)
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

					structName := strings.ToLower(typeSpec.Name.Name)
					fieldMap := make(map[string]string)

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}

					for _, field := range structType.Fields.List {
						fieldName := ""
						fieldType := ""

						for _, fieldNameIdent := range field.Names {
							if fieldNameIdent.Name != "ID" {
								fieldName = strings.ToLower(string(fieldNameIdent.Name[0])) + fieldNameIdent.Name[1:]
							} else {
								fieldName = strings.ToLower(fieldNameIdent.Name)
							}
						}

						switch fieldTypeExpr := field.Type.(type) {
						case *ast.Ident:
							fieldType = fieldTypeExpr.Name
						case *ast.SelectorExpr:
							fieldType = getFullSelectorExpr(fieldTypeExpr)
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

func getFullSelectorExpr(expr *ast.SelectorExpr) string {
	var pkgName string
	var identName string

	switch pkg := expr.X.(type) {
	case *ast.Ident:
		pkgName = pkg.Name
	case *ast.SelectorExpr:
		pkgName = getFullSelectorExpr(pkg)
	}

	identName = expr.Sel.Name

	return pkgName + "." + identName
}
