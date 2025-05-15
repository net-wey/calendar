package project

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
)

type ProjectRequestGetAll struct {
	ID uint `json:"id"`
}
type ProjectResponseGetAll struct {
	Status  string         `json:"status"`
	Error   string         `json:"error,omitempty"`
	Project entity.Project `json:"project,omitempty"`
}

type ProjectGetterGetAll interface {
	GetProject() (entity.Project, error)
}

func NewGetAllReportHandler(getter ProjectGetterGetAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ProjectResponseGetAll{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		project, err := getter.GetProject()
		if err != nil {
			if errors.Is(err, er.ErrReportNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ProjectResponseGetAll{
					Status: "error",
					Error:  "invalid request_id format",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ProjectResponseGetAll{
				Status: "error",
				Error:  "failed to get report",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ProjectResponseGetAll{
			Status:  "ok",
			Project: project,
		})
	}
}
