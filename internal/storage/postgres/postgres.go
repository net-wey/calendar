package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"time"

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

/////TASKS//////

func (s *Storage) TaskExists(uid uuid.UUID) (bool, error) {
	const op = "storage.postgres.TaskExists"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1)", uid).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: check task existence: %w", op, err)
	}
	return exists, nil
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

	exists, err := s.TaskExists(uid)
	if err != nil {
		return entity.Task{}, fmt.Errorf("%s: %w", op, err)
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

/////////DEVELOPERS/////////////

func (s *Storage) SaveDeveloper(developer entity.Developer) (uuid.UUID, error) {
	const op = "storage.postgres.SaveDeveloper"

	if developer.Name == "" || developer.LastName == "" {
		return uuid.Nil, fmt.Errorf("%s: %w", op, er.ErrInvalidDeveloperData)
	}

	uid := uuid.New()

	stmt, err := s.db.Prepare(
		`INSERT INTO developers(
			id,
			name,
			last_name,
			created_at,			
			modified_at,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		uid,
		developer.Name,
		developer.LastName,
		time.Now(),
	).Scan(&developer.CreatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return uid, nil
}

func (s *Storage) GetDeveloper(uid uuid.UUID) (entity.Developer, error) {
	const op = "storage.postgres.GetDeveloper"

	stmt, err := s.db.Prepare(`
		SELECT id, name, last_name, created_at
		FROM developers
		WHERE id = $1`)
	if err != nil {
		return entity.Developer{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var developer entity.Developer
	err = stmt.QueryRow(uid).Scan(
		&developer.ID,
		&developer.Name,
		&developer.LastName,
		&developer.CreatedAt,
	)
	if err != nil {
		return entity.Developer{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return developer, nil
}

func (s *Storage) GetDevelopers() ([]entity.Developer, error) {
	const op = "storage.postgres.GetDevelopers"

	stmt, err := s.db.Prepare(`
		SELECT id, name, last_name, created_at
		FROM developers`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var developers []entity.Developer
	for rows.Next() {
		var developer entity.Developer
		err := rows.Scan(
			&developer.ID,
			&developer.Name,
			&developer.LastName,
			&developer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		developers = append(developers, developer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return developers, nil
}

func (s *Storage) UpdateDeveloper(uid uuid.UUID, developer entity.Developer) error {
	const op = "storage.postgres.UpdateDeveloper"

	stmt, err := s.db.Prepare(
		`UPDATE developers SET 
		name = $1,
		last_name = $2
		WHERE id = $3`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		developer.Name,
		developer.LastName,
		uid,
	)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) SoftDeleteDeveloper(uid uuid.UUID) error {
	const op = "storage.postgres.SoftDeleteDeveloper"

	currentTime := time.Now()

	stmt, err := s.db.Prepare(`
        UPDATE developers 
        SET 
            deleted_at = $1,
            modified_at = $2 
        WHERE id = $3
        RETURNING modified_at`)

	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var updatedAt time.Time
	err = stmt.QueryRow("Удалён", currentTime, uid).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: developer not found: %w", op, err)
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

/////////////////////////////////
