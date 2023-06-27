package database

import (
	"fmt"
)

func InsertEntity(structName string, entity map[string]interface{}) {
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
	// query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT PRIMARY KEY, name VARCHAR(255))", tableName)

	// _, err := sqlConn.Exec(query)
	// if err != nil {
	// 	return fmt.Errorf("failed to create table: %w", err)
	// }

	fmt.Println(localStruct)

	return nil
}
