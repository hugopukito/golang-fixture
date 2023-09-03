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

type DatabaseParams struct {
	Name     string
	User     string
	Password string
	Ip       string
	Port     string
}

func InitDB(databaseParams *DatabaseParams) {

	if databaseParams.Name == "" {
		databaseParams.Name = "fixture"
	}
	if databaseParams.User == "" {
		databaseParams.User = "root"
	}
	if databaseParams.Ip == "" {
		databaseParams.Ip = "localhost"
	}
	if databaseParams.Port == "" {
		databaseParams.Port = "3306"
	}

	dataSource := databaseParams.User + ":" + databaseParams.Password + "@tcp(" + databaseParams.Ip + ":" + databaseParams.Port + ")/"

	sqlConn, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Panicln(color.Red + "Failed to connect to sql: " + err.Error() + color.Reset)
	}

	fmt.Println(color.Pink + "Drop database " + color.Yellow + databaseParams.Name + color.Pink + " if exists..." + color.Reset)
	fmt.Println(color.Pink + "Create database " + color.Yellow + databaseParams.Name + color.Pink + " if not exists..." + color.Reset)

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
		DatabaseName: databaseParams.Name,
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

	sqlConn, err = sql.Open("mysql", dataSource+databaseParams.Name)
	if err != nil {
		log.Panicln(color.Red + "Failed to connect to sql: " + err.Error() + color.Reset)
	}
}
