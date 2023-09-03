# Golang fixture data loader

## Description

This package will load yaml file, verify if each entity match with local structs of your project and insert rows in your database.

The entities in yaml files are named fixtures.

Steps of this package are :

- Create a database with the type structs defined for your entities in your project.
- Verify that all the fields for your yaml entities match with local structs defined in your project.
- Insert all the entities you described in yaml files

Working with mysql driver only (for the moment).

## Comming soon

- Permit {current} to be use as struct reference like dog_{current}

- Add * to multiply base on multiple entities the entity (like getting all entities starting with dog*, create entities for each one)

- Add primary key

- Add more drivers (now only mysql)

- Tests ahah

## Example

Create a main.go file

```go
package main

import (
	fixture "github.com/hugopukito/golang-fixture"
	"github.com/hugopukito/golang-fixture/database"
)

type Info struct {
  Text string
}

func main() {
	fixture.RunFixtures("my-fixtures", database.DatabaseParams{})
}
```

The 'RunFixtures' func takes two parameters:

- fixtureDirName: Give the location of your yaml files that describes entities.
- databaseParams: Choose the database name (that will be used ⚠️ or deleted if exist), user, password, ip, port.
- default values: dbName: 'fixture', user: 'root', password: '', ip: 'localhost', port: 3306


Do
```bash 
go mod init 'your-project-name'
```

Then
```bash
go get github.com/hugopukito/golang-fixture
```

Then create a folder that will contains entities in yaml (in the example) 'my-fixtures'

Create a yaml file name it as you want and insert entity that match struct from your project (Info struct in the example)

```yaml
Info:
  info1:
    text: "some text"
```

Run it
```bash 
go run main.go
```

## SQL types

This is the mapping used to create column types based on your project struct types.

So if you're using a type in Go that is not in this list, this will surely cause an error.

- `"int"`:       "INT",
- `"int8"`:      "TINYINT",
- `"int16"`:     "SMALLINT",
- `"int32"`:     "INT",
- `"int64"`:     "BIGINT",
- `"uint"`:      "INT UNSIGNED",
- `"uint8"`:     "TINYINT UNSIGNED",
- `"uint16"`:    "SMALLINT UNSIGNED",
- `"uint32"`:    "INT UNSIGNED",
- `"uint64"`:    "BIGINT UNSIGNED",
- `"float32"`:   "FLOAT",
- `"float64"`:   "DOUBLE",
- `"bool"`:      "BOOL",
- `"string"`:    "VARCHAR(255)",

Specials :

- `"time.Time"`: "TIMESTAMP",
- `"uuid.UUID"`: "VARCHAR(36)",

## Yaml

### Multiply

You can add range to create more entities without copy paste.

```yaml
Cat:
  cat{1..10}:
  	name: "chat{current}"
    color: "orange"
```

This will create 10 cats, the '{current}' will take the current number from the loop.

### New

You can use the key '{new{}}' for specials.

```yaml
dogNewDate:
  name: "new date dog"
  created: "{new{}}"
```

In this example it will generate a new time.

### Random

You can use random in three ways.

#### Range

For int and float with a range

```yaml
Cat:
  cat_random:
    name: "random_tester"
    color: "orange"
    tailLength: "{random{1..100}}"
```

#### List

For (normally) all with a list

```yaml
Cat:
  cat_random:
    name: "random_tester"
    color: "orange"
    tailLength: "{random{cat, cit, cot, cet}}"
```

#### Empty

For bool and specials (time and uuid)

```yaml
Cat:
  cat_random:
    name: "random_tester"
    color: "orange"
    tailLength: "{random{}}"
```

### Nested

You can make reference to another entity if you have a table referencing another one.

Just use the key '{ref{entityKey}}' with the entity key.

```yaml
Cat:
  cat_nested:
    name: "nested"
    info: "{ref{myInfo}}"
    more: "{ref{more1}}"

Info:
  myInfo:
    text: "some text"

More:
  more1:
    moreText: "more text"
```

It will then add '_id' to the column like 'info_id' and 'more_id' for the table cat.