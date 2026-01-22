package model

import "time"

type ExpenseSplit struct {
	ID        int       `json:"id"`
	ExpenseID int       `json:"expense_id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ExpenseSplitRequest struct {
	ExpenseID int     `json:"expense_id" binding:"required"`
	UserID    int     `json:"user_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

type ExpenseSplitResponse struct {
	ID        int     `json:"id"`
	ExpenseID int     `json:"expense_id"`
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
}
