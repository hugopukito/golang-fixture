# Golang fixture data loader

## Goal

Create a database with the type structs defined for your entities in your project.
Verify that all the fields for your yaml entities match with local structs defined in your project.
Insert all the entities you described in yaml files

## Description

This package will load yaml file, verify if each entity match with local structs of your project and insert rows in your database.

The entities in yaml files are named fixture.

## Example

Create a main.go file

```go
package main

import (
	fixture "github.com/hugopukito/golang-fixture"
)

type Info struct {
	Text string
}

func main() {
	fixture.RunFixtures("db-fixture-test", "my-fixtures")
}
```

Do ```bash 
go mod init 'your-project-name'
```
Then ```bash
go get github.com/hugopukito/golang-fixture
```

Then create a folder that will contains entities in yaml (in the example) 'my-fixtures'

Create a yaml file name it as you want and insert entity that match struct from your project (Info struct in the example)

```yaml
Info:
  info1:
    text: "some text"
```

Then do ```bash 
go run main.go
```
