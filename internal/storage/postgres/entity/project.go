package entity

import "time"

type Project struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	ModifiedAt  time.Time
	DeletedAt   *time.Time
}
