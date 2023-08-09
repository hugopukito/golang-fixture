package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/hugopukito/golang-fixture/color"
	"github.com/hugopukito/golang-fixture/database"
	"github.com/hugopukito/golang-fixture/funcs"

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
	err := loadConfig()
	if err != nil {
		fmt.Println(color.Red + "error loading .env file: " + err.Error() + color.Reset)
	}

	fmt.Println(color.Pink + "Testing connection to your sql..." + color.Reset)
	database.InitDB(config.databaseName)

	fmt.Println(color.Blue + "\nParsing your local structs..." + color.Reset)
	funcs.InitLocalStructs(config.structsPackageName)

	yamlFixtures, err := funcs.GetYamlStructs(config.fixtureDirName)
	if err != nil {
		fmt.Println(color.Red + "GetYamlStructs err: " + err.Error() + color.Reset)
		return
	}

	fmt.Println(color.Purple + "Parsing your fixtures... \n" + color.Reset)
	for _, yamlFixture := range yamlFixtures {
		funcs.ParseFixture(yamlFixture, config.databaseName)
	}
}

func loadConfig() error {
	err := godotenv.Load(".fixture.env")
	if err != nil {
		return err
	}

	config = &Config{
		databaseName:       getEnv("database_name"),
		structsPackageName: getEnv("structs_package_name"),
		fixtureDirName:     getEnv("fixture_dir_name"),
	}

	err = checkMissingParams(config)
	if err != nil {
		return err
	}
	return nil
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return ""
	}
	return value
}

func checkMissingParams(cfg *Config) error {
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
		return errors.New("error missing config parameter(s): " + strings.Join(missingParams, ", "))
	}
	return nil
}
