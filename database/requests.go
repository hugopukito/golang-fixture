package database

import (
	"errors"
	"fixture/color"
	"fmt"
	"log"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var specialTypes map[string]any
var firstColumns = []string{"id", "uid", "uuid"}
var currentRegex *regexp.Regexp
var randCommasRegex *regexp.Regexp
var randRangeRegex *regexp.Regexp

func init() {
	specialTypes = make(map[string]any)
	addFuncsToSpecialTypes()
	compileRegex()
}

func addFuncsToSpecialTypes() {
	generateUUID := func() any {
		return uuid.New()
	}
	specialTypes["uuid.UUID"] = generateUUID

	generateTimeStamp := func() any {
		return time.Now()
	}
	specialTypes["time.Time"] = generateTimeStamp
}

func compileRegex() {
	currentRegex, err = regexp.Compile(`\{current\}`)
	if err != nil {
		log.Println("Failed to compile regular expression:", err)
	}

	randCommasRegex, err = regexp.Compile(`\{random\{([^}]*)\}\}`)
	if err != nil {
		log.Println("Failed to compile regular expression:", err)
	}

	randRangeRegex, err = regexp.Compile(`\{random\{((?:\d+(?:\.\d+)?\.\.\d+(?:\.\d+)?))\}\}`)
	if err != nil {
		log.Println("Failed to compile regular expression:", err)
	}
}

func InsertEntity(structName string, entity map[string]any, localStruct map[string]string, occurrence int) error {
	columns := make([]string, 0)
	values := make([]any, 0)
	placeholders := make([]string, 0)

	// Special types auto generated if present in localStruct but not in entity
	for k, v := range localStruct {
		if _, exist := entity[k]; !exist {
			if value, ok := specialTypes[v]; ok {
				columns = append(columns, k)
				values = append(values, value.(func() any)())
				placeholders = append(placeholders, "?")
			}
		}
	}

	// Fields from entity
	for column, value := range entity {
		columns = append(columns, column)
		_, isString := value.(string)
		if isString {
			if currentRegex.MatchString(value.(string)) {
				values = append(values, strings.ReplaceAll(value.(string), "{current}", strconv.Itoa(occurrence)))
			} else if randCommasMatches := randCommasRegex.FindAllStringSubmatch(value.(string), -1); len(randCommasMatches) > 0 {
				if randRangeMatches := randRangeRegex.FindAllStringSubmatch(value.(string), -1); len(randRangeMatches) > 0 {
					innerContent := randCommasMatches[0][1]
					randomValues := splitAndTrim(innerContent, "..")
					randomVal, err := generateRandom(randomValues, localStruct[column])
					if err != nil {
						return err
					}
					values = append(values, randomVal)
				} else if strings.Contains(randCommasMatches[0][1], ",") {
					innerContent := randCommasMatches[0][1]
					randomValues := splitAndTrim(innerContent, ",")
					//castValuesInGoodType(randomValues, localStruct[column])
					values = append(values, getRandomElement(randomValues))
				} else {
					randomVal, err := generateRandom(nil, localStruct[column])
					if err != nil {
						return err
					}
					values = append(values, randomVal)
				}
			} else {
				values = append(values, value)
			}
		} else {
			values = append(values, value)
		}
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", structName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err = sqlConn.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func CheckTableExist(tableName string, dbName string) (bool, error) {
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"

	var count int
	err := sqlConn.QueryRow(query, dbName, tableName).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateTable(tableName string, localStruct map[string]string) error {
	fmt.Println(color.Purple+"Creating table ->", color.Yellow, tableName+"..."+color.Reset)
	columns := make([]string, 0, len(localStruct))
	idColumn := ""

	for columnName, columnType := range localStruct {
		sqlType, exist := GoSQLTypeMap[columnType]
		if !exist {
			return errors.New(color.Red + "sql type for type: " + color.Orange + columnType + color.Red + " doesn't exist")
		}
		// Just putting this on first column, doesn't change anything
		if containsString(columnName) {
			idColumn = fmt.Sprintf("%s %s", columnName, sqlType)
		} else {
			columns = append(columns, fmt.Sprintf("%s %s", columnName, sqlType))
		}
	}

	allColumns := append([]string{idColumn}, columns...)

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(allColumns, ", "))

	_, err := sqlConn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func containsString(target string) bool {
	for _, element := range firstColumns {
		if element == target {
			return true
		}
	}
	return false
}

func splitAndTrim(input, sep string) []string {
	var result []string
	parts := strings.Split(input, sep)
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getRandomElement(strings []string) string {
	if len(strings) == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(strings))
	randomElement := strings[randomIndex]

	return randomElement
}

func generateRandom(values []string, targetType string) (any, error) {
	if targetType == "bool" {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(2) == 1, nil
	} else if targetType == "int" || targetType == "float64" {
		return generateRandomNumber(values, targetType)
	}

	return nil, errors.New("can't generate random for type: " + targetType)
}

func generateRandomNumber(values []string, targetType string) (any, error) {
	rand.Seed(time.Now().UnixNano())

	if targetType == "int" {
		first, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, err
		}
		second, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}
		return rand.Intn(second) + first, nil
	} else if targetType == "float64" {
		first, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			return nil, err
		}
		second, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return nil, err
		}
		return math.Round((rand.Float64()*(second-first)+first)*100) / 100, nil
	}

	return nil, errors.New("problem with your random{x..y}")
}

// func castValuesInGoodType(randomValues []string, targetType string) ([]any, error) {
// 	typeChecker := make(map[string]func(string) (any, error))

// 	typeChecker["string"] = func(obj string) (any, error) {
// 		return obj, nil
// 	}
// 	typeChecker["int"] = func(obj string) (any, error) {
// 		return strconv.Atoi(obj)
// 	}
// 	typeChecker["float64"] = func(obj string) (any, error) {
// 		return strconv.ParseFloat(obj, 64)
// 	}
// 	typeChecker["bool"] = func(obj string) (any, error) {
// 		return strconv.ParseBool(obj)
// 	}

// 	var values []any

// 	for _, randomVal := range randomValues {
// 		if _func, ok := typeChecker[targetType]; ok {
// 			val, err := _func(randomVal)
// 			if err != nil {
// 				return nil, err
// 			} else {
// 				values = append(values, val)
// 			}
// 		}
// 	}

// 	return values, nil
// }
