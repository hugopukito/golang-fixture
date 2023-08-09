package fixture

import (
	"fmt"

	"github.com/hugopukito/golang-fixture/color"
	"github.com/hugopukito/golang-fixture/database"
	"github.com/hugopukito/golang-fixture/funcs"
)

// func main() {
// 	RunFixtures("db-fixture", "structs", "fixtures")
// }

func RunFixtures(databaseName, structsPackageName, fixtureDirName string) {

	fmt.Println(color.Pink + "Testing connection to your sql..." + color.Reset)
	database.InitDB(databaseName)

	fmt.Println(color.Blue + "\nParsing your local structs..." + color.Reset)
	funcs.InitLocalStructs(structsPackageName)

	yamlFixtures, err := funcs.GetYamlStructs(fixtureDirName)
	if err != nil {
		fmt.Println(color.Red + "GetYamlStructs err: " + err.Error() + color.Reset)
		return
	}

	fmt.Println(color.Purple + "Parsing your fixtures... \n" + color.Reset)
	for _, yamlFixture := range yamlFixtures {
		funcs.ParseFixture(yamlFixture, databaseName)
	}
}
