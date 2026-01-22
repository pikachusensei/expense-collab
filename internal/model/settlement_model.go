package model

import "time"

// Settlement represents a payment between two users to settle expenses
type Settlement struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	FromUserID  int       `json:"from_user_id"`
	ToUserID    int       `json:"to_user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SettlementRequest is the request body for creating a settlement
type SettlementRequest struct {
	GroupID     int     `json:"group_id"`
	FromUserID  int     `json:"from_user_id"`
	ToUserID    int     `json:"to_user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// SettlementResponse is the response body for settlement operations
type SettlementResponse struct {
	ID           int       `json:"id"`
	GroupID      int       `json:"group_id"`
	FromUserID   int       `json:"from_user_id"`
	FromUserName string    `json:"from_user_name"`
	ToUserID     int       `json:"to_user_id"`
	ToUserName   string    `json:"to_user_name"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

// GroupSettlementResponse shows all settlements in a group
type GroupSettlementResponse struct {
	Settlements []SettlementResponse `json:"settlements"`
	TotalAmount float64              `json:"total_amount"`
}
