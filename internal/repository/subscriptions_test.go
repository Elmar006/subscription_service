package repository

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Elmar006/subscription_service/internal/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var testRepo SubscriptionRepository

func TestMain(m *testing.M) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/subscription_test?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	testRepo = NewSubscriptionRepo(db)
	os.Exit(m.Run())
}

func createTestSubscription(t *testing.T) *model.Subscription {
	sub := &model.Subscription{
		ServiceName: "Test Service",
		Price:       555,
		UserID:      uuid.New().String(),
		StartDate:   "2026-01-01",
		EndDate:     "2026-01-31",
		CreatedAt:   time.Now(),
	}

	if err := testRepo.Create(sub); err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	return sub
}

func TestCreateSubscription(t *testing.T) {
	sub := createTestSubscription(t)

	check, err := testRepo.GetByID(sub.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if check.ID != sub.ID {
		t.Errorf("Expected ID %v, got %v", sub.ID, check.ID)
	}
}

func TestUpdateSubscription(t *testing.T) {
	sub := createTestSubscription(t)
	sub.Price = 777
	sub.ServiceName = "Music TestService"

	if err := testRepo.Update(sub); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	check, err := testRepo.GetByID(sub.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}

	if check.Price != sub.Price || check.ServiceName != sub.ServiceName {
		t.Error("Update did not persist changes")
	}
}

func TestDeleteSubscription(t *testing.T) {
	sub := createTestSubscription(t)

	if err := testRepo.Delete(sub.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	check, _ := testRepo.GetByID(sub.ID)
	if check != nil {
		t.Errorf("Subscription was not deleted")
	}
}

func TestListByUser(t *testing.T) {
	userID := uuid.New().String()

	for i := 0; i < 3; i++ {
		sub := &model.Subscription{
			ServiceName: "Service " + strconv.Itoa(i+1),
			Price:       100 + i*10,
			UserID:      userID,
			StartDate:   "2026-01-01",
			EndDate:     "2026-01-31",
			CreatedAt:   time.Now(),
		}
		if err := testRepo.Create(sub); err != nil {
			t.Fatalf("Failed to create subscription %d: %v", i, err)
		}
	}

	subs, err := testRepo.ListByUser(userID)
	if err != nil {
		t.Fatalf("ListByUser failed: %v", err)
	}

	if len(subs) != 3 {
		t.Errorf("Expected 3 subscriptions, got %d", len(subs))
	}
}

func TestTotal(t *testing.T) {
	userID := uuid.New().String()
	serviceName := "ServiceTotalTest"

	prices := []int{100, 200, 300}
	for _, p := range prices {
		sub := &model.Subscription{
			ServiceName: serviceName,
			Price:       p,
			UserID:      userID,
			StartDate:   "2026-01-01",
			EndDate:     "2026-12-31",
			CreatedAt:   time.Now(),
		}
		if err := testRepo.Create(sub); err != nil {
			t.Fatalf("Failed to create subscription for total: %v", err)
		}
	}

	from, _ := time.Parse("2006-01-02", "2026-01-01")
	to, _ := time.Parse("2006-01-02", "2026-12-31")
	total, err := testRepo.Total(&userID, &serviceName, from, to)
	if err != nil {
		t.Fatalf("Total calculation failed: %v", err)
	}

	expected := 0
	for _, p := range prices {
		expected += p
	}

	if total != expected {
		t.Errorf("Expected total %d, got %d", expected, total)
	}
}
