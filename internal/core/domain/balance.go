package domain

import "time"

type Balance struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
