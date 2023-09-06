# TODO

## Cool

Permit to have a text type column instead of varchar(255) to generate bigger strings

## Bug

When adding new field in entity with nested fields, doesn't work first time but working other times.
 -> Probably because was using CompanyId with uuid.UUID type instead of Company with Company type
 -> Didn't got it since a while

## Features

- Permit {current} to be use as struct reference like dog_{current}

- Add * to multiply base on multiple entities the entity (like getting all entities starting with dog*, create entities for each one)

## Later if necessary

- Add primary key

- Add more drivers (now only mysql)

- Tests ahah