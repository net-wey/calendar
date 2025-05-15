package project

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
)

type ProjectRequestPost struct {
	ID uint `json:"id"`
}
type ProjectResponsePost struct {
	Status  string         `json:"status"`
	Error   string         `json:"error,omitempty"`
	Project entity.Project `json:"project,omitempty"`
}

type ProjectSaverPost interface {
	SaveProject(project entity.Project) error
}

func NewProjectHandler(saver ProjectSaverPost) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ProjectResponsePost{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		var req ProjectRequestPost
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponsePost{
				Status: "error",
				Error:  "failed to decode request",
			})
			return
		}

		if req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ProjectResponsePost{
				Status: "error",
				Error:  "id is required",
			})
			return
		}

		project := entity.Project{
			ID: req.ID,
		}

		err := saver.SaveProject(project)
		if err != nil {
			if errors.Is(err, er.ErrInvalidDeveloperData) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ProjectResponsePost{
					Status: "error",
					Error:  "invalid project data",
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProjectResponsePost{
				Status: "error",
				Error:  "failed to save project",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ProjectResponsePost{
			Status:  "ok",
			Project: project,
		})
	}
}
