package project

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProjectUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
}

// ProjectResponse - структура ответа для проектов
type ProjectResponse struct {
	Status    string         `json:"status"`
	Error     string         `json:"error,omitempty"`
	Project   entity.Project `json:"project,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

type ProjectUpdater interface {
	GetProjectByID(ID uint) (entity.Project, error)
	UpdateProject(ID uint, project entity.Project) error
}

func NewUpdateProjectHandler(updater ProjectUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем метод запроса
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		// Извлекаем ID проекта из URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "invalid URL format",
			})
			return
		}

		projectID, err := strconv.ParseUint(pathParts[2], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "invalid project ID format",
			})
			return
		}

		// Парсим тело запроса
		var req ProjectUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "invalid request body",
			})
			return
		}

		// Валидация входных данных
		if len(req.Name) < 2 || len(req.Name) > 100 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "name must be between 2 and 100 characters",
			})
			return
		}

		if len(req.Description) > 500 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "description must not exceed 500 characters",
			})
			return
		}

		// Проверяем существование проекта
		existingProject, err := updater.GetProjectByID(uint(projectID))
		if err != nil {
			if errors.Is(err, er.ErrProjectNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ProjectResponse{
					Status: "error",
					Error:  "project not found",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "failed to get project",
			})
			return
		}

		// Обновляем проект
		updatedProject := entity.Project{
			ID:          uint(projectID),
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   existingProject.CreatedAt,
		}

		if err := updater.UpdateProject(uint(projectID), updatedProject); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProjectResponse{
				Status: "error",
				Error:  "failed to update project",
			})
			return
		}

		// Возвращаем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ProjectResponse{
			Status:    "success",
			Project:   updatedProject,
			Timestamp: time.Now(),
		})
	}
}
