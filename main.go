package main

import (
	"fixture/color"
	"fixture/database"
	"fixture/funcs"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	databaseName       string
	structsPackageName string
	fixtureDirName     string
}

var config *Config

func main() {

	fmt.Println(color.Blue + "Parsing your config.conf file...\n" + color.Reset)
	loadConfig()

	fmt.Println(color.Pink + "Testing connection to your sql..." + color.Reset)
	database.InitDB(config.databaseName)

	fmt.Println(color.Blue + "\nParsing your local structs..." + color.Reset)
	funcs.InitLocalStructs(config.structsPackageName)

	yamlFixtures, err := funcs.GetYamlStructs(config.fixtureDirName)
	if err != nil {
		log.Panicln(color.Red + "GetYamlStructs err: " + err.Error() + color.Reset)
	}

	fmt.Println(color.Purple + "Parsing your fixtures... \n" + color.Reset)
	for _, yamlFixture := range yamlFixtures {
		funcs.ParseFixture(yamlFixture, config.databaseName)
	}
}

func loadConfig() {
	err := godotenv.Load(".fixture.env")
	if err != nil {
		log.Panicln(color.Red + "error loading .env file: " + err.Error() + color.Reset)
	}

	config = &Config{
		databaseName:       getEnv("database_name"),
		structsPackageName: getEnv("structs_package_name"),
		fixtureDirName:     getEnv("fixture_dir_name"),
	}

	checkMissingParams(config)
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return ""
	}
	return value
}

func checkMissingParams(cfg *Config) {
	missingParams := []string{}
	configType := reflect.TypeOf(*cfg)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := reflect.ValueOf(cfg).Elem().FieldByName(field.Name).String()

		if fieldValue == "" {
			missingParams = append(missingParams, field.Name)
		}
	}

	if len(missingParams) > 0 {
		log.Panicln(color.Red + "error missing config parameter(s): " + strings.Join(missingParams, ", ") + color.Reset)
	}
}
