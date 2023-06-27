package structs

import (
	"time"

	"github.com/google/uuid"
)

type Dog struct {
	ID    uuid.UUID
	Name  string
	Breed string
	Size  int
	Date  time.Time
}

type Cat struct {
	ID    uuid.UUID
	Name  string
	Color string
}
