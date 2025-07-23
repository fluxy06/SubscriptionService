package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sbs/models"
	"sbs/services"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// SubscriptionHandler отвечает за обработку HTTP-запросов для подписок
type SubscriptionHandler struct {
	service *services.SubscriptionService
}

// Конструктор
func NewSubscriptionHandler(service *services.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// --- Вспомогательная структура для ответа с форматированными датами ---
type subscriptionResponse struct {
	ID          int     `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func toSubscriptionResponse(s *models.Subscription) subscriptionResponse {
	var endDateStr *string
	if s.EndDate != nil {
		str := models.FormatMonthYear(*s.EndDate)
		endDateStr = &str
	}
	return subscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   models.FormatMonthYear(s.StartDate),
		EndDate:     endDateStr,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
}

// POST /subscriptions
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат данных: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем валидность user_id как UUID
	if _, err := uuid.Parse(req.UserID); err != nil {
		http.Error(w, "user_id должен быть валидным UUID", http.StatusBadRequest)
		return
	}

	startDate, err := models.ParseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, "неверный формат start_date: "+err.Error(), http.StatusBadRequest)
		return
	}

	var endDatePtr *time.Time
	if req.EndDate != nil {
		endDate, err := models.ParseMonthYear(*req.EndDate)
		if err != nil {
			http.Error(w, "неверный формат end_date: "+err.Error(), http.StatusBadRequest)
			return
		}
		endDatePtr = &endDate
	}

	sub := models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}

	if err := h.service.Create(&sub); err != nil {
		log.Printf("Ошибка создания подписки: %v\n", err)
		http.Error(w, "ошибка при создании: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := toSubscriptionResponse(&sub)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GET /subscriptions/{id}
func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "ошибка при получении: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if subscription == nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	resp := toSubscriptionResponse(subscription)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PUT /subscriptions/{id}
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	// Декодируем тело запроса
	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат данных: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем валидность user_id сразу после декодирования
	if _, err := uuid.Parse(req.UserID); err != nil {
		http.Error(w, "user_id должен быть валидным UUID", http.StatusBadRequest)
		return
	}

	// Парсим даты start_date и end_date в time.Time
	startDate, err := models.ParseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, "неверный формат start_date: "+err.Error(), http.StatusBadRequest)
		return
	}

	var endDatePtr *time.Time
	if req.EndDate != nil {
		endDate, err := models.ParseMonthYear(*req.EndDate)
		if err != nil {
			http.Error(w, "неверный формат end_date: "+err.Error(), http.StatusBadRequest)
			return
		}
		endDatePtr = &endDate
	}

	// Формируем объект подписки для обновления
	sub := models.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}

	// Выполняем обновление через сервис
	if err := h.service.Update(&sub); err != nil {
		http.Error(w, "ошибка при обновлении: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// После успешного обновления получаем актуальную запись из базы
	updatedSub, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "ошибка при получении обновлённой записи: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedSub == nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	// Формируем ответ с актуальными created_at и updated_at
	resp := toSubscriptionResponse(updatedSub)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DELETE /subscriptions/{id}
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "ошибка при удалении: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /subscriptions
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	subs, err := h.service.List()
	if err != nil {
		http.Error(w, "ошибка при получении списка: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make([]subscriptionResponse, 0, len(subs))
	for i := range subs {
		resp = append(resp, toSubscriptionResponse(&subs[i]))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /subscriptions/sum?user_id=...&service_name=...&start=MM-YYYY&end=MM-YYYY
func (h *SubscriptionHandler) Sum(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDStr := query.Get("user_id")
	serviceName := query.Get("service_name")
	startStr := query.Get("start")
	endStr := query.Get("end")

	if userIDStr == "" || serviceName == "" || startStr == "" || endStr == "" {
		http.Error(w, "missing required query parameters", http.StatusBadRequest)
		return
	}

	// Проверяем user_id на валидность UUID
	if _, err := uuid.Parse(userIDStr); err != nil {
		http.Error(w, "invalid user_id format, must be UUID", http.StatusBadRequest)
		return
	}

	// Парсим даты в формате MM-YYYY
	startDate, err := time.Parse("01-2006", startStr)
	if err != nil {
		http.Error(w, "invalid start date format, must be MM-YYYY", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("01-2006", endStr)
	if err != nil {
		http.Error(w, "invalid end date format, must be MM-YYYY", http.StatusBadRequest)
		return
	}

	sum, err := h.service.SumSubscriptions(userIDStr, serviceName, startDate, endDate)
	if err != nil {
		http.Error(w, "error calculating sum: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"sum": sum}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
