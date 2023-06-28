package database

import (
	"errors"
	"fixture/color"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func InsertEntity(structName string, entity map[string]interface{}, keysWithTypeUUID []string) error {
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)

	for _, v := range keysWithTypeUUID {
		columns = append(columns, v)
		values = append(values, uuid.New())
		placeholders = append(placeholders, "?")
	}

	for column, value := range entity {
		columns = append(columns, column)
		values = append(values, value)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		structName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	_, err := sqlConn.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func CheckTableExist(tableName string, dbName string) (bool, error) {
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"

	var count int
	err := sqlConn.QueryRow(query, dbName, "your_table_name").Scan(&count)
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
		sqlType, exist := goSQLTypeMap[columnType]
		if !exist {
			return errors.New(color.Red + "sql type for type: " + color.Orange + columnType + color.Red + " doesn't exist")
		}
		if columnName == "id" {
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
