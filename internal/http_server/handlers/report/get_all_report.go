package report

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
)

type ReportResponseGetAll struct {
	Status  string          `json:"status"`
	Error   string          `json:"error,omitempty"`
	Reports []entity.Report `json:"reports,omitempty"`
}

type ReportGetterGetAll interface {
	GetReport() ([]entity.Report, error)
}

func NewGetAllReportHandler(getter ReportGetterGetAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ReportResponseGetAll{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		reports, err := getter.GetReport()
		if err != nil {
			if errors.Is(err, er.ErrReportNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ReportResponseGetAll{
					Status: "error",
					Error:  "reports not found",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ReportResponseGetAll{
				Status: "error",
				Error:  "failed to get reports",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ReportResponseGetAll{
			Status:  "ok",
			Reports: reports,
		})
	}
}
