package model

type Balance struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id"`
	GroupID int     `json:"group_id"`
	Amount  float64 `json:"amount"`
}

type BalanceRequest struct {
	UserID  int `json:"user_id" binding:"required"`
	GroupID int `json:"group_id" binding:"required"`
}

type BalanceResponse struct {
	UserID  int     `json:"user_id"`
	GroupID int     `json:"group_id"`
	Amount  float64 `json:"amount"`
}

type UserBalanceResponse struct {
	UserID   int     `json:"user_id"`
	UserName string  `json:"user_name"`
	Amount   float64 `json:"amount"`
}
