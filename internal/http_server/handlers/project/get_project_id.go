package project

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
	"strconv"
	"strings"
)

// ReportResponseGet - структура ответа для получения отчета
type ProjectResponseGet struct {
	Status  string         `json:"status"`
	Error   string         `json:"error,omitempty"`
	Project entity.Project `json:"project,omitempty"`
}

// ReportGetter - интерфейс для получения отчета
type ProjectGetter interface {
	GetProjectById(id uint) (entity.Project, error)
}

// NewGetReportByIdHandler создает обработчик для получения отчета по ID
func NewGetProjectByIdHandler(getter ProjectGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем метод запроса
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ProjectResponseGet{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		// Извлекаем ID отчета из URL
		projectIDStr := strings.Split(r.URL.Path, "/")
		if len(projectIDStr) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponseGet{
				Status: "error",
				Error:  "project ID is required",
			})
			return
		}

		// Парсим UUID
		projectID, err := strconv.ParseUint(projectIDStr[2], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponseGet{
				Status: "error",
				Error:  "invalid project ID format",
			})
			return
		}

		// Получаем отчет из хранилища
		project, err := getter.GetProjectById(uint(projectID))
		if err != nil {
			if errors.Is(err, er.ErrReportNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ProjectResponseGet{
					Status: "error",
					Error:  "project not found",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProjectResponseGet{
				Status: "error",
				Error:  "failed to get project",
			})
			return
		}

		// Возвращаем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ProjectResponseGet{
			Status:  "success",
			Project: project,
		})
	}
}
