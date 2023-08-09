package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/hugopukito/golang-fixture/color"

	_ "github.com/go-sql-driver/mysql"
)

var sqlConn *sql.DB
var err error

func InitDB(dbName string) {

	sqlConn, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/")
	if err != nil {
		log.Panicln(color.Red + "Failed to connect to sql: " + err.Error() + color.Reset)
	}

	fmt.Println(color.Pink + "drop database " + color.Yellow + dbName + color.Pink + " if exists..." + color.Reset)
	fmt.Println(color.Pink + "create database " + color.Yellow + dbName + color.Pink + " if not exists..." + color.Reset)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Failed to get caller information")
		return
	}

	scriptFile := filepath.Join(filepath.Dir(filename), "/reset.sql")

	scriptContent, err := os.ReadFile(scriptFile)
	if err != nil {
		log.Panicln(color.Red + "Failed to read the SQL script file: " + err.Error() + color.Reset)
	}

	// Create a template from the script content
	tmpl, err := template.New("script").Parse(string(scriptContent))
	if err != nil {
		log.Panicln(color.Red + "Failed to parse the SQL script template: " + err.Error() + color.Reset)
	}

	// Prepare the data to substitute in the template
	data := struct {
		DatabaseName string
	}{
		DatabaseName: dbName,
	}

	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		log.Panicln(color.Red + "Failed to execute the SQL script template: " + err.Error() + color.Reset)
	}

	// Split the script into separate SQL statements
	sqlStatements := strings.Split(output.String(), ";")

	// Execute each SQL statement separately
	for _, stmt := range sqlStatements {
		trimmedStmt := strings.TrimSpace(stmt)
		if trimmedStmt != "" {
			_, err = sqlConn.Exec(trimmedStmt)
			if err != nil {
				log.Panicln(color.Red + "Failed to execute SQL statement: " + err.Error() + color.Reset)
			}
		}
	}

	sqlConn, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/"+dbName)
	if err != nil {
		log.Panicln(color.Red + "Failed to connect to sql: " + err.Error() + color.Reset)
	}
}
