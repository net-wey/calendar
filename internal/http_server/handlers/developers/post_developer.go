package handlers

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type DeveloperRequest struct {
	Name      string     `json:"name"`
	LastName  string     `json:"last_name"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type DeveloperResponse struct {
	Status      string    `json:"status"`
	Error       string    `json:"error,omitempty"`
	DeveloperID uuid.UUID `json:"developer_id,omitempty"`
}

type DeveloperSaver interface {
	SaveDeveloper(developer entity.Developer) (uuid.UUID, error)
}

func NewDeveloperHandler(saver DeveloperSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		var req DeveloperRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "failed to decode request",
			})
			return
		}

		if req.Name == "" || req.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "firstname and last_name are required",
			})
			return
		}

		developer := entity.Developer{
			Name:      req.Name,
			LastName:  req.LastName,
			DeletedAt: req.DeletedAt,
		}

		developerID, err := saver.SaveDeveloper(developer)
		if err != nil {
			if errors.Is(err, er.ErrInvalidDeveloperData) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(DeveloperResponse{
					Status: "error",
					Error:  "invalid developer data",
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "failed to save developer",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(DeveloperResponse{
			Status:      "ok",
			DeveloperID: developerID,
		})
	}
}
