package repositorypg

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
)

type ExpenseRepositoryPG struct {
	DB *sql.DB
}

func NewExpenseRepositoryPG(db *sql.DB) *ExpenseRepositoryPG {
	return &ExpenseRepositoryPG{DB: db}
}

func (r *ExpenseRepositoryPG) CreateExpense(expense *model.Expense) (*model.Expense, error) {
	query := `
		INSERT INTO expenses (group_id, paid_by_id, amount, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, group_id, paid_by_id, amount, description, created_at, updated_at
	`

	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		expense.GroupID,
		expense.PaidByID,
		expense.Amount,
		expense.Description,
		expense.CreatedAt,
		expense.UpdatedAt,
	).Scan(&expense.ID, &expense.GroupID, &expense.PaidByID, &expense.Amount, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)

	if err != nil {
		log.Printf("Error creating expense: %v", err)
		return nil, err
	}

	return expense, nil
}

func (r *ExpenseRepositoryPG) GetExpenseByID(id int) (*model.Expense, error) {
	query := `
		SELECT id, group_id, paid_by_id, amount, description, created_at, updated_at
		FROM expenses
		WHERE id = $1
	`

	expense := &model.Expense{}
	err := r.DB.QueryRow(query, id).Scan(
		&expense.ID,
		&expense.GroupID,
		&expense.PaidByID,
		&expense.Amount,
		&expense.Description,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expense not found")
		}
		log.Printf("Error getting expense by ID: %v", err)
		return nil, err
	}

	return expense, nil
}

func (r *ExpenseRepositoryPG) GetExpensesByGroupID(groupID int) ([]*model.Expense, error) {
	query := `
		SELECT id, group_id, paid_by_id, amount, description, created_at, updated_at
		FROM expenses
		WHERE group_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		log.Printf("Error getting expenses by group ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []*model.Expense
	for rows.Next() {
		expense := &model.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.GroupID,
			&expense.PaidByID,
			&expense.Amount,
			&expense.Description,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning expense: %v", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating expenses: %v", err)
		return nil, err
	}

	return expenses, nil
}

func (r *ExpenseRepositoryPG) GetExpensesByUserID(userID int) ([]*model.Expense, error) {
	query := `
		SELECT id, group_id, paid_by_id, amount, description, created_at, updated_at
		FROM expenses
		WHERE paid_by_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error getting expenses by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []*model.Expense
	for rows.Next() {
		expense := &model.Expense{}
		err := rows.Scan(
			&expense.ID,
			&expense.GroupID,
			&expense.PaidByID,
			&expense.Amount,
			&expense.Description,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning expense: %v", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating expenses: %v", err)
		return nil, err
	}

	return expenses, nil
}

func (r *ExpenseRepositoryPG) UpdateExpense(expense *model.Expense) (*model.Expense, error) {
	query := `
		UPDATE expenses
		SET amount = $1, description = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, group_id, paid_by_id, amount, description, created_at, updated_at
	`

	expense.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		expense.Amount,
		expense.Description,
		expense.UpdatedAt,
		expense.ID,
	).Scan(&expense.ID, &expense.GroupID, &expense.PaidByID, &expense.Amount, &expense.Description, &expense.CreatedAt, &expense.UpdatedAt)

	if err != nil {
		log.Printf("Error updating expense: %v", err)
		return nil, err
	}

	return expense, nil
}

func (r *ExpenseRepositoryPG) DeleteExpense(id int) error {
	query := `DELETE FROM expenses WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting expense: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("expense not found")
	}

	return nil
}
