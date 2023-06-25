package structs

import "github.com/google/uuid"

type Dog struct {
	ID    uuid.UUID
	Name  string
	Breed string
}

type Cat struct {
	ID    uuid.UUID
	Name  string
	Breed string
}
