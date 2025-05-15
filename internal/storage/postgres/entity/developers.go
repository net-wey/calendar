package entity

import (
	"time"

	"github.com/google/uuid"
)

type Developer struct {
	ID         uuid.UUID
	Name       string
	LastName   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	DeletedAt  *time.Time
}
