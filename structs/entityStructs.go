package structs

import "github.com/google/uuid"

type Dog struct {
	ID    uuid.UUID
	Name  string
	Breed string
	Size  int
}

type Cat struct {
	ID    uuid.UUID
	Name  string
	Color string
}
