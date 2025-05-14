package postgres

import "errors"

var (
	// ErrDeveloperNotFound returns when developer not found in storage
	ErrDeveloperNotFound = errors.New("developer not found")

	// ErrTaskNotFound returns when task not found in storage
	ErrTaskNotFound = errors.New("task not found")

	// ErrProjectNotFound returns when project not found in storage
	ErrProjectNotFound = errors.New("project not found")

	// ErrReportNotFound returns when report not found in storage
	ErrReportNotFound = errors.New("report not found")

	// ErrInvalidDeveloperData returns when developer data is invalid
	ErrInvalidDeveloperData = errors.New("invalid developer data")

	// ErrInvalidTaskData returns when task data is invalid
	ErrInvalidTaskData = errors.New("invalid task data")

	// ErrInvalidProjectData returns when project data is invalid
	ErrInvalidProjectData = errors.New("invalid project data")

	// ErrInvalidReportData returns when report data is invalid
	ErrInvalidReportData = errors.New("invalid report data")
)
