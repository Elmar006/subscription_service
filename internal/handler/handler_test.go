package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Elmar006/subscription_service/internal/config"
	"github.com/Elmar006/subscription_service/internal/db"
	"github.com/Elmar006/subscription_service/internal/model"
	"github.com/Elmar006/subscription_service/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func setupHandler(t *testing.T) (*SubscriptionHandler, *model.Subscription, repository.SubscriptionRepository) {
	cfg := &config.Config{
		DBHost: "localhost",
		DBPort: "5432",
		DBUser: "postgres",
		DBPass: "postgres",
		DBName: "subscription_test",
	}
	database := db.Connect(cfg)
	repo := repository.NewSubscriptionRepo(database)
	h := NewSubscriptionHandler(repo)

	userID := uuid.New().String()

	sub := &model.Subscription{
		ServiceName: "Test Service",
		Price:       888,
		UserID:      userID,
		StartDate:   "2026-01-01",
		EndDate:     "2026-01-31",
		CreatedAt:   time.Now(),
	}
	if err := repo.Create(sub); err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	return h, sub, repo
}

func TestCreateSub(t *testing.T) {
	h, _, _ := setupHandler(t)

	body := map[string]interface{}{
		"service_name": "Music Plus",
		"price":        500,
		"user_id":      uuid.New().String(),
		"start_date":   "2026-02-01",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(data))
	w := httptest.NewRecorder()

	h.CreateSubscription(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", w.Code)
	}

	var sub model.Subscription
	if err := json.NewDecoder(w.Body).Decode(&sub); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if sub.ServiceName != "Music Plus" || sub.Price != 500 {
		t.Errorf("Response mismatch, got %v", sub)
	}
	if sub.ID == "" {
		t.Errorf("Expected generated ID, got empty")
	}
}

func TestGetByIDSub(t *testing.T) {
	h, sub, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+sub.ID, nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", sub.ID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	h.GetByIDSubscription(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var got model.Subscription
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if got.ID != sub.ID {
		t.Errorf("Expected ID %v, got %v", sub.ID, got.ID)
	}
}

func TestListByUser(t *testing.T) {
	h, sub, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions?user_id="+sub.UserID, nil)
	w := httptest.NewRecorder()

	h.GetSubscription(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var subs []model.Subscription
	if err := json.NewDecoder(w.Body).Decode(&subs); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	if len(subs) == 0 || subs[0].ID != sub.ID {
		t.Errorf("Expected subscription list to contain %+v, got %+v", sub.ID, subs)
	}
}

func TestUpdateSubscription(t *testing.T) {
	h, sub, _ := setupHandler(t)

	update := map[string]interface{}{
		"service_name": "Updated Service",
		"price":        999,
	}
	data, _ := json.Marshal(update)
	req := httptest.NewRequest(http.MethodPut, "/subscriptions/"+sub.ID, bytes.NewReader(data))
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", sub.ID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	h.UpdateByIDSubscription(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var updated model.Subscription
	if err := json.NewDecoder(w.Body).Decode(&updated); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	if updated.ServiceName != "Updated Service" || updated.Price != 999 {
		t.Errorf("Update failed, got %+v", updated)
	}
}

func TestDeleteSubscription(t *testing.T) {
	h, sub, _ := setupHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+sub.ID, nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", sub.ID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	h.DeleteSubscription(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected 204 No Content, got %d", w.Code)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/subscriptions/"+sub.ID, nil)
	wGet := httptest.NewRecorder()
	reqGet = reqGet.WithContext(ctx)
	h.GetByIDSubscription(wGet, reqGet)

	if wGet.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 Not Found after delete, got %d", wGet.Code)
	}
}

func TestGetSubscriptionTotal(t *testing.T) {
	h, sub, repo := setupHandler(t)

	sub2 := &model.Subscription{
		ServiceName: "Test Service",
		Price:       200,
		UserID:      sub.UserID,
		StartDate:   "2026-01-15",
		EndDate:     "2026-02-15",
		CreatedAt:   time.Now(),
	}
	repo.Create(sub2)

	req := httptest.NewRequest(http.MethodGet, "/subscriptions/total?user_id="+sub.UserID+"&from=2026-01-01&to=2026-01-31", nil)
	w := httptest.NewRecorder()

	h.GetSubscriptionTotal(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	var resp map[string]int
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := sub.Price + sub2.Price
	if resp["total"] != expected {
		t.Errorf("Expected total %d, got %d", expected, resp["total"])
	}
}
