package repository

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/Elmar006/subscription_service/internal/model"
	"github.com/Elmar006/subscription_service/logger"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id string) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id string) error
	ListByUser(userID string) ([]*model.Subscription, error)
	Total(userID *string, serviceName *string, from, to time.Time) (int, error)
}

type subscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (s *subscriptionRepo) Create(sub *model.Subscription) error {
	startDate, err := parseDate(sub.StartDate)
	if err != nil {
		return err
	}

	var endDate sql.NullTime
	if sub.EndDate != "" {
		t, err := parseDate(sub.EndDate)
		if err != nil {
			return err
		}
		endDate = sql.NullTime{Time: t, Valid: true}
	}

	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}

	_, err = s.db.Exec(
		`INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		sub.ID, sub.ServiceName, sub.Price, sub.UserID, startDate, endDate, sub.CreatedAt,
	)
	if err != nil {
		logger.L().Errorf("Error inserting subscription: %v", err)
		return err
	}

	return nil
}

func (s *subscriptionRepo) GetByID(id string) (*model.Subscription, error) {
	row := s.db.QueryRow(
		`SELECT id, service_name, price, user_id, start_date, end_date, created_at
		 FROM subscriptions WHERE id = $1`, id)

	sub := &model.Subscription{}
	var startDate time.Time
	var endDate sql.NullTime

	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &startDate, &endDate, &sub.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.L().Errorf("Error fetching subscription: %v", err)
		return nil, err
	}

	sub.StartDate = startDate.Format("2006-01-02")
	if endDate.Valid {
		sub.EndDate = endDate.Time.Format("2006-01-02")
	}

	return sub, nil
}

func (s *subscriptionRepo) Update(sub *model.Subscription) error {
	startDate, err := parseDate(sub.StartDate)
	if err != nil {
		return err
	}

	var endDate sql.NullTime
	if sub.EndDate != "" {
		t, err := parseDate(sub.EndDate)
		if err != nil {
			return err
		}
		endDate = sql.NullTime{Time: t, Valid: true}
	}

	_, err = s.db.Exec(
		`UPDATE subscriptions SET service_name=$1, price=$2, start_date=$3, end_date=$4 WHERE id=$5`,
		sub.ServiceName, sub.Price, startDate, endDate, sub.ID,
	)
	if err != nil {
		logger.L().Errorf("Error updating subscription: %v", err)
		return err
	}

	return nil
}

func (s *subscriptionRepo) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM subscriptions WHERE id=$1`, id)
	if err != nil {
		logger.L().Errorf("Error deleting subscription: %v", err)
		return err
	}
	return nil
}

func (s *subscriptionRepo) ListByUser(userID string) ([]*model.Subscription, error) {
	rows, err := s.db.Query(
		`SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions WHERE user_id=$1`,
		userID)
	if err != nil {
		logger.L().Errorf("Error listing subscriptions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		var startDate time.Time
		var endDate sql.NullTime
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &startDate, &endDate, &sub.CreatedAt); err != nil {
			return nil, err
		}
		sub.StartDate = startDate.Format("2006-01-02")
		if endDate.Valid {
			sub.EndDate = endDate.Time.Format("2006-01-02")
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (s *subscriptionRepo) Total(userID *string, serviceName *string, from, to time.Time) (int, error) {
	query := `SELECT COALESCE(SUM(price),0) FROM subscriptions WHERE start_date >= $1 AND start_date <= $2`
	args := []interface{}{from, to}
	argIndex := 3

	if userID != nil {
		query += " AND user_id = $" + strconv.Itoa(argIndex)
		args = append(args, *userID)
		argIndex++
	}

	if serviceName != nil {
		query += " AND service_name = $" + strconv.Itoa(argIndex)
		args = append(args, *serviceName)
		argIndex++
	}

	var total int
	err := s.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		logger.L().Errorf("Error calculating total: %v", err)
		return 0, err
	}

	return total, nil
}
