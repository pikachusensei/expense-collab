package model

import "time"

type GroupMember struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupMemberRequest struct {
	GroupID int    `json:"group_id" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

type GroupMemberResponse struct {
	ID       int       `json:"id"`
	GroupID  int       `json:"group_id"`
	UserID   int       `json:"user_id"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	AddedAt  time.Time `json:"added_at"`
}
