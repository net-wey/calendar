package report

import (
	"encoding/json"
	"errors"
	er "goproject/internal/storage"
	"goproject/internal/storage/postgres/entity"
	"net/http"
)

type ReportRequestPost struct {
	ID uint `json:"id" validate:"required"`
}

type ReportResponsePost struct {
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	ReportID uint   `json:"report_id,omitempty"`
}

type ReportSaverPost interface {
	SaveReport(report entity.Report) error
}

func NewReportHandler(saver ReportSaverPost) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ReportResponsePost{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		var req ReportRequestPost
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ReportResponsePost{
				Status: "error",
				Error:  "failed to decode request",
			})
			return
		}

		if req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ReportResponsePost{
				Status: "error",
				Error:  "developer_id is required",
			})
			return
		}

		report := entity.Report{
			ID: req.ID,
		}

		err := saver.SaveReport(report)
		if err != nil {
			if errors.Is(err, er.ErrInvalidDeveloperData) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ReportResponsePost{
					Status: "error",
					Error:  "invalid report data",
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ReportResponsePost{
				Status: "error",
				Error:  "failed to save report",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ReportResponsePost{
			Status:   "ok",
			ReportID: report.ID,
		})
	}
}
