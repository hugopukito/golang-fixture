package fixture

import (
	"fmt"

	"github.com/hugopukito/golang-fixture/color"
	"github.com/hugopukito/golang-fixture/database"
	"github.com/hugopukito/golang-fixture/funcs"
)

// Create (or use) a database with the type structs defined for your entities in your project.
//
// Verify that all the fields for your yaml entities match with local structs defined in your project.
//
// Insert all the entities you described in yaml files.
//
// Parameters:
//
// - databaseName: Choose the database name that will be used ⚠️ or deleted if exist
//
// - fixtureDirName: Give the location of your yaml files that describes entities
func RunFixtures(databaseName, fixtureDirName string) {

	fmt.Println(color.Pink + "Testing connection to your sql..." + color.Reset)
	database.InitDB(databaseName)

	fmt.Println(color.Blue + "\nParsing your local structs..." + color.Reset)
	funcs.InitLocalStructs()

	yamlFixtures, err := funcs.GetYamlStructs(fixtureDirName)
	if err != nil {
		fmt.Println(color.Red + "GetYamlStructs err: " + err.Error() + color.Reset)
		return
	}

	fmt.Println(color.Purple + "Parsing your fixtures... \n" + color.Reset)
	funcs.ParseFixture(yamlFixtures, databaseName)
}
