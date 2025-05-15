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

func (s *Storage) SaveTask(task entity.Task) (int, error) {
	const op = "storage.postgres.SaveTask"

	if task.Name == "" || task.EstimatePlaned <= 0 {
		return 0, fmt.Errorf("%s: %w", op, er.ErrInvalidTaskData)
	}

	stmt, err := s.db.Prepare(
		`INSERT INTO tasks(
            report_id,
            project_id,
            name,
            developer_note,
            estimate_planed,
            estimate_progress,
            start_timestamp,
            end_timestamp
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at`) // Changed to return both id and created_at
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(
		task.ReportID,
		task.ProjectID,
		task.Name,
		task.DeveloperNote,
		task.EstimatePlaned,
		task.EstimateProgress,
		task.StartTimestamp,
		task.EndTimestamp,
	).Scan(&id, &task.CreatedAt) // Scan both id and created_at
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetTaskByID(ID uint) (entity.Task, error) {
	const op = "storage.postgres.GetTask"

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
	err = stmt.QueryRow(ID).Scan(
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

func (s *Storage) GetTasks() ([]entity.Task, error) {
	const op = "storage.postgres.GetTasks"

	stmt, err := s.db.Prepare(`
        SELECT id, report_id, project_id, name, developer_note, 
               estimate_planed, estimate_progress, 
               start_timestamp, end_timestamp, created_at
        FROM tasks
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
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
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return tasks, nil
}

func (s *Storage) GetTasksByReportID(ID uint) ([]entity.Task, error) {
	const op = "storage.postgres.GetTasksByReportID"

	stmt, err := s.db.Prepare(`
		SELECT id, report_id, project_id, name, developer_note, 
               estimate_planed, estimate_progress, 
               start_timestamp, end_timestamp, created_at
        FROM tasks
		WHERE report_id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
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
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
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

func (s *Storage) GetDeveloperByID(uid uuid.UUID) (entity.Developer, error) {
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

/////////////////////////////////REPORTS//////////////////////////////

func (s *Storage) SaveReport(report entity.Report) error {
	const op = "storage.postgres.SaveReport"

	stmt, err := s.db.Prepare(
		`INSERT INTO reports(
    		developer_id,
    	c	reated_at
    	) VALUES ($1, $2)
    	RETURNING id, created_at`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		report.DeveloperID,
		time.Now(),
	).Scan(&report.ID, &report.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}

func (s *Storage) GetReport() ([]entity.Report, error) {
	const op = "storage.postgres.GetReport"

	stmt, err := s.db.Prepare(`
    SELECT id, developer_id, created_at
    FROM reports
    ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var reports []entity.Report
	for rows.Next() {
		var report entity.Report
		err := rows.Scan(
			&report.ID,
			&report.DeveloperID,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return reports, nil
}

func (s *Storage) GetReportById(id uint) (entity.Report, error) {
	const op = "storage.postgres.GetReportById"

	stmt, err := s.db.Prepare(`
    SELECT id, developer_id, created_at
    FROM reports  
    WHERE id = $1`)
	if err != nil {
		return entity.Report{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var report entity.Report
	err = stmt.QueryRow(id).Scan(
		&report.ID,
		&report.DeveloperID,
		&report.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Report{}, fmt.Errorf("%s: %w", op, er.ErrReportNotFound)
		}
		return entity.Report{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return report, nil
}

func (s *Storage) GetReportsByDeveloperID(developerID uuid.UUID) ([]entity.Report, error) {
	const op = "storage.postgres.GetReportsByDeveloperID"

	stmt, err := s.db.Prepare(`
        SELECT id, developer_id, created_at
        FROM reports
        WHERE developer_id = $1
        ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(developerID)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var reports []entity.Report
	for rows.Next() {
		var report entity.Report
		err := rows.Scan(
			&report.ID,
			&report.DeveloperID,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return reports, nil
}

/////////////////////////////////PROJECTS//////////////////////////////

func (s *Storage) SaveProject(project entity.Project) error {
	const op = "storage.postgres.SaveProject"

	stmt, err := s.db.Prepare(`
    INSERT INTO projects(
      name,
      description,
      created_at
    ) VALUES ($1, $2, $3)
    RETURNING id, created_at`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		project.Name,
		project.Description,
		time.Now(),
	).Scan(&project.ID, &project.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}

func (s *Storage) GetProject() ([]entity.Project, error) {
	const op = "storage.postgres.GetProject"

	stmt, err := s.db.Prepare(`
    SELECT id, name, description, created_at
    FROM projects
    ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	defer rows.Close()

	var projects []entity.Project
	for rows.Next() {
		var project entity.Project
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return projects, nil
}

func (s *Storage) GetProjectByID(ID uint) (entity.Project, error) {
	const op = "storage.postgres.GetProjectByID"

	stmt, err := s.db.Prepare(`
    SELECT id, name, description, created_at
    FROM projects
    WHERE id = $1`)
	if err != nil {
		return entity.Project{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var project entity.Project
	err = stmt.QueryRow(ID).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Project{}, fmt.Errorf("%s: %w", op, er.ErrProjectNotFound)
		}
		return entity.Project{}, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return project, nil
}

func (s *Storage) UpdateProject(ID uint, project entity.Project) error {
	const op = "storage.postgres.UpdateProject"

	stmt, err := s.db.Prepare(`
        UPDATE projects 
        SET name = $1, description = $2, modified_at = NOW()
        WHERE id = $3`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(project.Name, project.Description, ID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
