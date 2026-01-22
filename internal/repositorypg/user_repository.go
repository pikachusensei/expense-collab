package repositorypg

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/lib/pq"
	"github.com/shreyansh/expense-go-collab-backend/internal/model"
)

type UserRepositoryPG struct {
	DB *sql.DB
}

func NewUserRepositoryPG(db *sql.DB) *UserRepositoryPG {
	return &UserRepositoryPG{DB: db}
}

func (r *UserRepositoryPG) CreateUser(user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, created_at, updated_at
	`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		user.Email,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, fmt.Errorf("user with email already exists")
		}
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPG) GetUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPG) GetUserByID(id int) (*model.User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPG) GetAllUsers() ([]*model.User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating users: %v", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryPG) UpdateUser(user *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET name = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, email, name, created_at, updated_at
	`

	user.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		user.Name,
		user.UpdatedAt,
		user.ID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPG) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
