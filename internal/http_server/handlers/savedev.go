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
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"last_name"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// DeveloperResponse представляет структуру ответа на запрос создания разработчика
type DeveloperResponse struct {
	Status      string    `json:"status"`
	Error       string    `json:"error,omitempty"`
	DeveloperID uuid.UUID `json:"developer_id,omitempty"`
}

// DeveloperSaver определяет интерфейс для сохранения информации о разработчике
type DeveloperSaver interface {
	SaveDeveloper(developer entity.Developer) (uuid.UUID, error)
}

// NewDeveloperHandler создает новый обработчик HTTP для сохранения разработчика
func NewDeveloperHandler(saver DeveloperSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем метод запроса
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "method not allowed",
			})
			return
		}

		// Декодируем тело запроса
		var req DeveloperRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "failed to decode request",
			})
			return
		}

		// Валидация данных
		if req.Firstname == "" || req.Lastname == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(DeveloperResponse{
				Status: "error",
				Error:  "firstname and last_name are required",
			})
			return
		}

		// Подготавливаем сущность Developer
		developer := entity.Developer{
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			DeletedAt: req.DeletedAt,
		}

		// Сохраняем разработчика
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

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(DeveloperResponse{
			Status:      "ok",
			DeveloperID: developerID,
		})
	}
}
