package repositorypg

import (
	"database/sql"
	"log"
	"time"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
)

type ExpenseSplitRepositoryPG struct {
	DB *sql.DB
}

func NewExpenseSplitRepositoryPG(db *sql.DB) *ExpenseSplitRepositoryPG {
	return &ExpenseSplitRepositoryPG{DB: db}
}

func (r *ExpenseSplitRepositoryPG) CreateSplit(split *model.ExpenseSplit) (*model.ExpenseSplit, error) {
	query := `
		INSERT INTO expense_splits (expense_id, user_id, amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, expense_id, user_id, amount, created_at, updated_at
	`

	split.CreatedAt = time.Now()
	split.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		split.ExpenseID,
		split.UserID,
		split.Amount,
		split.CreatedAt,
		split.UpdatedAt,
	).Scan(&split.ID, &split.ExpenseID, &split.UserID, &split.Amount, &split.CreatedAt, &split.UpdatedAt)

	if err != nil {
		log.Printf("Error creating split: %v", err)
		return nil, err
	}

	return split, nil
}

func (r *ExpenseSplitRepositoryPG) GetSplitsByExpenseID(expenseID int) ([]*model.ExpenseSplit, error) {
	query := `
		SELECT id, expense_id, user_id, amount, created_at, updated_at
		FROM expense_splits
		WHERE expense_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, expenseID)
	if err != nil {
		log.Printf("Error getting splits by expense ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var splits []*model.ExpenseSplit
	for rows.Next() {
		split := &model.ExpenseSplit{}
		err := rows.Scan(
			&split.ID,
			&split.ExpenseID,
			&split.UserID,
			&split.Amount,
			&split.CreatedAt,
			&split.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning split: %v", err)
			return nil, err
		}
		splits = append(splits, split)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating splits: %v", err)
		return nil, err
	}

	return splits, nil
}

func (r *ExpenseSplitRepositoryPG) GetSplitsByUserID(userID int) ([]*model.ExpenseSplit, error) {
	query := `
		SELECT id, expense_id, user_id, amount, created_at, updated_at
		FROM expense_splits
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error getting splits by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var splits []*model.ExpenseSplit
	for rows.Next() {
		split := &model.ExpenseSplit{}
		err := rows.Scan(
			&split.ID,
			&split.ExpenseID,
			&split.UserID,
			&split.Amount,
			&split.CreatedAt,
			&split.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning split: %v", err)
			return nil, err
		}
		splits = append(splits, split)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating splits: %v", err)
		return nil, err
	}

	return splits, nil
}

func (r *ExpenseSplitRepositoryPG) DeleteSplitsByExpenseID(expenseID int) error {
	query := `DELETE FROM expense_splits WHERE expense_id = $1`

	_, err := r.DB.Exec(query, expenseID)
	if err != nil {
		log.Printf("Error deleting splits: %v", err)
		return err
	}

	return nil
}

func (r *ExpenseSplitRepositoryPG) UpdateSplit(split *model.ExpenseSplit) (*model.ExpenseSplit, error) {
	query := `
		UPDATE expense_splits
		SET amount = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, expense_id, user_id, amount, created_at, updated_at
	`

	split.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		split.Amount,
		split.UpdatedAt,
		split.ID,
	).Scan(&split.ID, &split.ExpenseID, &split.UserID, &split.Amount, &split.CreatedAt, &split.UpdatedAt)

	if err != nil {
		log.Printf("Error updating split: %v", err)
		return nil, err
	}

	return split, nil
}
