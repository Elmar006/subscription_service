package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Elmar006/subscription_service/internal/model"
	"github.com/Elmar006/subscription_service/internal/repository"
	"github.com/Elmar006/subscription_service/logger"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

// @Summary Create a new subscription
// @Description Create a new subscription with the input payload
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Subscription data"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
func (s *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub model.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if sub.ServiceName == "" || sub.Price < 0 || sub.UserID == "" || sub.StartDate == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sub.CreatedAt = time.Now()

	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}

	if err := s.repo.Create(&sub); err != nil {
		http.Error(w, "Failed to create subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// @Summary Get subscription by ID
// @Description Get subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (s *SubscriptionHandler) GetByIDSubscription(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "Invalid id parameter is required", http.StatusBadRequest)
		return
	}

	sub, err := s.repo.GetByID(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if sub == nil {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// GetSubscription godoc
// @Summary List all subscriptions for a user
// @Description Returns all subscriptions belonging to a specific user
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param user_id query string true "User ID (UUID)"
// @Success 200 {array} model.Subscription
// @Failure 400 {string} string "Invalid user_id parameter is required"
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions [get]
func (s *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Invalid user_id parameter is required", http.StatusBadRequest)
		return
	}

	sub, err := s.repo.ListByUser(userIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// @Summary Update subscription by ID
// @Description Update subscription fields
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body model.Subscription true "Subscription data"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (s *SubscriptionHandler) UpdateByIDSubscription(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	existing, err := s.repo.GetByID(idParam)
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

// GetSubscriptionTotal godoc
// @Summary Get total price of subscriptions
// @Description Returns the total sum of subscriptions for a user with optional filters
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name filter"
// @Param from query string false "Start date filter (YYYY-MM-DD)"
// @Param to query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} map[string]int "Returns total as JSON {\"total\":123}"
// @Failure 400 {string} string "Invalid date format"
// @Failure 500 {string} string "Internal server error"
// @Router /subscriptions/total [get]
func (s *SubscriptionHandler) GetSubscriptionTotal(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	userIDStr := q.Get("user_id")
	serviceNameStr := q.Get("service_name")
	fromStr := q.Get("from")
	toStr := q.Get("to")

	var userIDPtr *string
	if userIDStr != "" {
		userIDPtr = &userIDStr
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
	json.NewEncoder(w).Encode(map[string]int{"total": total})
}

// @Summary Delete subscription by ID
// @Description Delete subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (s *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	subCheck, err := s.repo.GetByID(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if subCheck == nil {
		logger.L().Error("The subscription you want to delete does not exist")
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	if err := s.repo.Delete(idParam); err != nil {
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
