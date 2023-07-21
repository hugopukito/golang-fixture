package database

import (
	"errors"
	"fixture/color"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var specialTypes map[string]any
var firstColumns = []string{"id", "uid", "uuid"}

func init() {
	specialTypes = make(map[string]any)
	addFuncsToSpecialTypes()
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

	currentRegex, err := regexp.Compile(`\{current\}`)
	if err != nil {
		log.Println("Failed to compile regular expression:", err)
	}

	// Fields from entity
	for column, value := range entity {
		columns = append(columns, column)
		_, isString := value.(string)
		if isString {
			if currentRegex.MatchString(value.(string)) {
				values = append(values, strings.ReplaceAll(value.(string), "{current}", strconv.Itoa(occurrence)))
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
