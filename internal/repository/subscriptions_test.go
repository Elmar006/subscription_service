package repository

import (
	"database/sql"
	"log"
	"os"
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
		UserID:      uuid.New(),
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
