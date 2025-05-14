package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID               uuid.UUID
	ReportID         uuid.UUID
	ProjectID        uuid.UUID
	Name             string
	DeveloperNote    string
	EstimatePlaned   int
	EstimateProgress int
	StartTimestamp   time.Time
	EndTimestamp     time.Time
	CreatedAt        time.Time
}
