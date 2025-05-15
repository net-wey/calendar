package entity

import (
	"time"
)

type Task struct {
	ID               uint
	ReportID         uint
	ProjectID        uint
	Name             string
	DeveloperNote    string
	EstimatePlaned   int
	EstimateProgress int
	StartTimestamp   time.Time
	EndTimestamp     time.Time
	CreatedAt        time.Time
}
