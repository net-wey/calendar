package entity

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID          uint
	DeveloperID uuid.UUID
	CreatedAt   time.Time
	ModifiedAt  time.Time
	DeletedAt   *time.Time
}
