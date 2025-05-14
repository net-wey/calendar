package postgres

import (
	"database/sql"
	"fmt"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(url string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetCalendar(uid uuid.UUID) error {
	const op = "storage.postgres.GetCalendar"

	// Сначала проверяем существование задачи
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", uid).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: check task existence: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: %w", op, er.ErrTaskNotFound)
	}

	return nil
}

func (s *Storage) SaveTask(task entity.Task) (uuid.UUID, error) {
	const op = "storage.postgres.SaveTask"

	if task.Name == "" || task.EstimatePlaned <= 0 {
		return uuid.Nil, fmt.Errorf("%s: %w", op, er.ErrInvalidTaskData)
	}

	stmt, err := s.db.Prepare(
		`INSERT INTO tasks(
			id,
			report_id,
			project_id,
			name,
			developer_note,
			estimate_planed,
			estimate_progress,
			start_timestamp,
			end_timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at`)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	uid := uuid.New()
	err = stmt.QueryRow(
		uid,
		task.ReportID,
		task.ProjectID,
		task.Name,
		task.DeveloperNote,
		task.EstimatePlaned,
		task.EstimateProgress,
		task.StartTimestamp,
		task.EndTimestamp,
	).Scan(&task.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return uid, nil
}

func (s *Storage) GetTask(uid uuid.UUID) (entity.Task, error) {
	const op = "storage.postgres.GetTask"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", uid).Scan(&exists)
	if err != nil {
		return entity.Task{}, fmt.Errorf("%s: check task existence: %w", op, err)
	}
	if !exists {
		return entity.Task{}, fmt.Errorf("%s: %w", op, er.ErrTaskNotFound)
	}

	stmt, err := s.db.Prepare(`
		SELECT id, report_id, project_id, name, developer_note, 
			   estimate_planed, estimate_progress, 
			   start_timestamp, end_timestamp, created_at
		FROM tasks 
		WHERE id = $1`)
	if err != nil {
		return entity.Task{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var task entity.Task
	err = stmt.QueryRow(uid).Scan(
		&task.ID,
		&task.ReportID,
		&task.ProjectID,
		&task.Name,
		&task.DeveloperNote,
		&task.EstimatePlaned,
		&task.EstimateProgress,
		&task.StartTimestamp,
		&task.EndTimestamp,
		&task.CreatedAt,
	)
	if err != nil {
		return entity.Task{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return task, nil
}
