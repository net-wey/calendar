package report

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
type ReportResponseGet struct {
	Status string        `json:"status"`
	Error  string        `json:"error,omitempty"`
	Report entity.Report `json:"report,omitempty"`
}

// ReportGetter - интерфейс для получения отчета
type ReportGetter interface {
	GetReportById(id uint) (entity.Report, error)
}

// NewGetReportByIdHandler создает обработчик для получения отчета по ID
func NewGetReportByIdHandler(getter ReportGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем метод запроса
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ReportResponseGet{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		// Извлекаем ID отчета из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ReportResponseGet{
				Status: "error",
				Error:  "report ID is required",
			})
			return
		}

		// Парсим uint
		reportID, err := strconv.ParseUint(parts[2], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ReportResponseGet{
				Status: "error",
				Error:  "invalid report ID format",
			})
			return
		}

		// Получаем отчет из хранилища
		report, err := getter.GetReportById(uint(reportID))
		if err != nil {
			if errors.Is(err, er.ErrReportNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(ReportResponseGet{
					Status: "error",
					Error:  "report not found",
				})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ReportResponseGet{
				Status: "error",
				Error:  "failed to get report",
			})
			return
		}

		// Возвращаем успешный ответ
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ReportResponseGet{
			Status: "success",
			Report: report,
		})
	}
}
