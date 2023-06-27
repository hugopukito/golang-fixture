package database

import (
	"errors"
	"fixture/color"
	"fmt"
	"strings"
)

func InsertEntity(structName string, entity map[string]interface{}) {
	// TODO insert
}

func CheckTableExist(tableName string) (bool, error) {
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = ?"

	var count int
	err := sqlConn.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if table exists: %w", err)
	}

	return count > 0, nil
}

func CreateTable(tableName string, localStruct map[string]string) error {
	columns := make([]string, 0, len(localStruct))

	for column, columnType := range localStruct {
		sqlType, exist := goSQLTypeMap[columnType]
		if !exist {
			return errors.New(color.Red + "sql type for type: " + color.Orange + columnType + color.Red + " doesn't exist")
		}
		columns = append(columns, fmt.Sprintf("%s %s", column, sqlType))
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))

	_, err := sqlConn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
