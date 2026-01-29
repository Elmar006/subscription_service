package model

import "time"

type Subscription struct {
	// @json id
	// @format uuid
	ID          string    `json:"id" example:"4658b3ad-0323-4d4d-854c-05da025bf9ef"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}
