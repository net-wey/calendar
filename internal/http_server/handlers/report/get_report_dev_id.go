package report

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type DevReportResponseGet struct {
	Status  string          `json:"status"`
	Error   string          `json:"error"`
	Reports []entity.Report `json:"reports"`
	Count   int             `json:"count"`
}

type DevReportsGetter interface {
	GetDeveloperById(uid uuid.UUID) (entity.Developer, error)
	GetReportsByDeveloperID(developerID uuid.UUID) ([]entity.Report, error)
}

func NewGetDeveloperReportsHandler(getter DevReportsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(DevReportResponseGet{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DevReportResponseGet{
				Status: "error",
				Error:  "invalid URL path",
			})
			return
		}

		developerIDStr := pathParts[2]
		developerID, err := uuid.Parse(developerIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DevReportResponseGet{
				Status: "error",
				Error:  "invalid URL path",
			})
			return
		}

		_, err = getter.GetDeveloperById(developerID)
		if err != nil {
			if errors.Is(err, er.ErrDeveloperNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(DevReportResponseGet{
					Status: "error",
					Error:  "developer not found",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(DevReportResponseGet{
				Status: "error",
				Error:  "failed to get developer",
			})
			return
		}

		reports, err := getter.GetReportsByDeveloperID(developerID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(DevReportResponseGet{
				Status: "error",
				Error:  "failed to get reports",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DevReportResponseGet{
			Status:  "success",
			Reports: reports,
			Count:   len(reports),
		})

	}
}
