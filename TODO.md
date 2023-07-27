# TODO

## Features

- Keep order when creating table (maybe look on localStruct)
 -> dans localStructsParser, faire un slice qui garde l'ordre des noms des tables (Cat, Dog) 
    et un sous slice pour l'ordre des colonnes (id, name)
- Implement nested structs
 -> give id (if not specified, than will be write in sql) for each struct than needs id before starting inserting entities
- Permit {current} to be use as struct reference like dog_{current}
- Add * to multiply base on multiple entities the entity (like getting all entities starting with dog*, create entities for each one)

## Later if necessary

- Add primary key