package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Elmar006/subscription_service/internal/model"
	"github.com/Elmar006/subscription_service/internal/repository"
	"github.com/Elmar006/subscription_service/logger"
)

type SubscriptionHandler struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

func (s *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if sub.ServiceName == "" || sub.Price < 0 || sub.UserID == uuid.Nil || sub.StartDate == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sub.CreatedAt = time.Now()

	if err := s.repo.Create(&sub); err != nil {
		http.Error(w, "Failed to create subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (s *SubscriptionHandler) GetByIDSubscription(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "Invalid id parameter is required", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub, err := s.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sub == nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

func (s *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Invalid user_id parameter is required", http.StatusBadRequest)
		return
	}
	usId, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	sub, err := s.repo.ListByUser(usId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

func (s *SubscriptionHandler) UpdateByIDSubscription(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	var sub model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if sub.ServiceName != "" {
		existing.ServiceName = sub.ServiceName
	}
	if sub.Price >= 0 {
		existing.Price = sub.Price
	}
	if sub.StartDate != "" {
		existing.StartDate = sub.StartDate
	}
	if sub.EndDate != "" {
		existing.EndDate = sub.EndDate
	}

	if err := s.repo.Update(existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.L().Info("Subscription updated successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

func (s *SubscriptionHandler) GetSubscriptionTotal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	userIDStr := q.Get("user_id")
	serviceNameStr := q.Get("service_name")
	fromStr := q.Get("from")
	toStr := q.Get("to")

	var userIDPtr *uuid.UUID
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}
		userIDPtr = &userID
	}

	var serviceNamePtr *string
	if serviceNameStr != "" {
		serviceNamePtr = &serviceNameStr
	}

	var from time.Time
	if fromStr != "" {
		t, err := parseDate(fromStr)
		if err != nil {
			http.Error(w, "Invalid 'from' date", http.StatusBadRequest)
			return
		}
		from = t
	} else {
		from = time.Time{}
	}

	var to time.Time
	if toStr != "" {
		t, err := parseDate(toStr)
		if err != nil {
			http.Error(w, "Invalid 'to' date", http.StatusBadRequest)
			return
		}
		to = t
	} else {
		to = time.Now()
	}

	total, err := s.repo.Total(userIDPtr, serviceNamePtr, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"total": total})
}

func (s *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subCheck, err := s.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if subCheck == nil {
		logger.L().Error("The subscription you want to delete does not exist")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	if err := s.repo.Delete(id); err != nil {
		http.Error(w, "Error deleting an entry", http.StatusInternalServerError)
		return
	}

	logger.L().Info("Subscription record successfully deleted")
	w.WriteHeader(http.StatusNoContent)
}

func parseDate(date string) (time.Time, error) {
	if len(date) == 7 {
		return time.Parse("2006-01", date)
	}

	return time.Parse("2006-01-02", date)
}
