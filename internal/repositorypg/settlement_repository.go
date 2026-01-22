package repositorypg

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
)

type SettlementRepositoryPG struct {
	DB *sql.DB
}

func NewSettlementRepositoryPG(db *sql.DB) *SettlementRepositoryPG {
	return &SettlementRepositoryPG{DB: db}
}

func (r *SettlementRepositoryPG) CreateSettlement(settlement *model.Settlement) (*model.Settlement, error) {
	query := `
		INSERT INTO settlements (group_id, from_user_id, to_user_id, amount, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, group_id, from_user_id, to_user_id, amount, description, created_at, updated_at
	`

	settlement.CreatedAt = time.Now()
	settlement.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		settlement.GroupID,
		settlement.FromUserID,
		settlement.ToUserID,
		settlement.Amount,
		settlement.Description,
		settlement.CreatedAt,
		settlement.UpdatedAt,
	).Scan(&settlement.ID, &settlement.GroupID, &settlement.FromUserID, &settlement.ToUserID,
		&settlement.Amount, &settlement.Description, &settlement.CreatedAt, &settlement.UpdatedAt)

	if err != nil {
		log.Printf("Error creating settlement: %v", err)
		return nil, fmt.Errorf("failed to create settlement: %v", err)
	}

	return settlement, nil
}

func (r *SettlementRepositoryPG) GetSettlementByID(id int) (*model.Settlement, error) {
	query := `
		SELECT id, group_id, from_user_id, to_user_id, amount, description, created_at, updated_at
		FROM settlements
		WHERE id = $1
	`

	settlement := &model.Settlement{}
	err := r.DB.QueryRow(query, id).Scan(
		&settlement.ID, &settlement.GroupID, &settlement.FromUserID, &settlement.ToUserID,
		&settlement.Amount, &settlement.Description, &settlement.CreatedAt, &settlement.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("settlement not found")
		}
		return nil, fmt.Errorf("failed to get settlement: %v", err)
	}

	return settlement, nil
}

func (r *SettlementRepositoryPG) GetSettlementsByGroupID(groupID int) ([]*model.Settlement, error) {
	query := `
		SELECT id, group_id, from_user_id, to_user_id, amount, description, created_at, updated_at
		FROM settlements
		WHERE group_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query settlements: %v", err)
	}
	defer rows.Close()

	var settlements []*model.Settlement
	for rows.Next() {
		settlement := &model.Settlement{}
		if err := rows.Scan(
			&settlement.ID, &settlement.GroupID, &settlement.FromUserID, &settlement.ToUserID,
			&settlement.Amount, &settlement.Description, &settlement.CreatedAt, &settlement.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan settlement: %v", err)
		}
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

func (r *SettlementRepositoryPG) GetSettlementsByUserID(userID int) ([]*model.Settlement, error) {
	query := `
		SELECT id, group_id, from_user_id, to_user_id, amount, description, created_at, updated_at
		FROM settlements
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query settlements: %v", err)
	}
	defer rows.Close()

	var settlements []*model.Settlement
	for rows.Next() {
		settlement := &model.Settlement{}
		if err := rows.Scan(
			&settlement.ID, &settlement.GroupID, &settlement.FromUserID, &settlement.ToUserID,
			&settlement.Amount, &settlement.Description, &settlement.CreatedAt, &settlement.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan settlement: %v", err)
		}
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

func (r *SettlementRepositoryPG) GetAllSettlements() ([]*model.Settlement, error) {
	query := `
		SELECT id, group_id, from_user_id, to_user_id, amount, description, created_at, updated_at
		FROM settlements
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query settlements: %v", err)
	}
	defer rows.Close()

	var settlements []*model.Settlement
	for rows.Next() {
		settlement := &model.Settlement{}
		if err := rows.Scan(
			&settlement.ID, &settlement.GroupID, &settlement.FromUserID, &settlement.ToUserID,
			&settlement.Amount, &settlement.Description, &settlement.CreatedAt, &settlement.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan settlement: %v", err)
		}
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}
