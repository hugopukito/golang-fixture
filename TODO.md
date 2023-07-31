# TODO

## Package

Add logic to package working
 -> Variables from env
 -> Work with path to get yaml and structs

## Features

- Implement nested structs
  -> Need to solve problem of adding 0 to entityName after single entities in yaml
  -> give id (if not specified, than will be write in sql) for each struct than needs id before starting inserting entities

- Permit {current} to be use as struct reference like dog_{current}

- Add * to multiply base on multiple entities the entity (like getting all entities starting with dog*, create entities for each one)

### Later if necessary

- Do loading bar when lots of entities that take more that 1s to create

- Add primary key