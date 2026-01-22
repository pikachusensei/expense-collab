package model

import "time"

type Expense struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	PaidByID    int       `json:"paid_by_id"`
	Amount      float64   `json:"amount" binding:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpenseRequest struct {
	GroupID     int     `json:"group_id" binding:"required"`
	PaidByID    int     `json:"paid_by_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}

type ExpenseResponse struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	PaidByID    int       `json:"paid_by_id"`
	PaidByName  string    `json:"paid_by_name"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
