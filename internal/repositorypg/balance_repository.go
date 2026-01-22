package repositorypg

import (
	"database/sql"
	"log"
)

type BalanceRepositoryPG struct {
	DB *sql.DB
}

func NewBalanceRepositoryPG(db *sql.DB) *BalanceRepositoryPG {
	return &BalanceRepositoryPG{DB: db}
}

func (r *BalanceRepositoryPG) GetBalance(userID, groupID int) (float64, error) {
	return r.GetUserBalanceInGroup(userID, groupID)
}

func (r *BalanceRepositoryPG) GetUserBalanceInGroup(userID, groupID int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(es.amount), 0) - COALESCE(SUM(CASE WHEN e.paid_by_id = $1 THEN e.amount ELSE 0 END), 0)
		FROM expenses e
		LEFT JOIN expense_splits es ON e.id = es.expense_id
		WHERE e.group_id = $2 AND es.user_id = $1
	`

	var balance float64
	err := r.DB.QueryRow(query, userID, groupID).Scan(&balance)
	if err != nil {
		log.Printf("Error getting balance: %v", err)
		return 0, err
	}

	return balance, nil
}

func (r *BalanceRepositoryPG) GetGroupBalances(groupID int) (map[int]float64, error) {
	query := `
		SELECT DISTINCT u.id,
			COALESCE(SUM(es.amount), 0) - COALESCE(SUM(CASE WHEN e.paid_by_id = u.id THEN e.amount ELSE 0 END), 0) as balance
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		LEFT JOIN expense_splits es ON u.id = es.user_id
		LEFT JOIN expenses e ON es.expense_id = e.id AND e.group_id = $1
		WHERE gm.group_id = $1
		GROUP BY u.id
	`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		log.Printf("Error getting group balances: %v", err)
		return nil, err
	}
	defer rows.Close()

	balances := make(map[int]float64)
	for rows.Next() {
		var userID int
		var balance float64
		err := rows.Scan(&userID, &balance)
		if err != nil {
			log.Printf("Error scanning balance: %v", err)
			return nil, err
		}
		balances[userID] = balance
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating balances: %v", err)
		return nil, err
	}

	return balances, nil
}

func (r *BalanceRepositoryPG) CalculateBalances(groupID int) error {
	// This would typically involve more complex calculations
	// For now, we're using the query-based approach above
	return nil
}
