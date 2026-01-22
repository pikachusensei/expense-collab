package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "expense_tracker"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connection successful")

	createTables(db)

	return db
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS groups (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			creator_id INTEGER NOT NULL REFERENCES users(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS group_members (
			id SERIAL PRIMARY KEY,
			group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(group_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			paid_by_id INTEGER NOT NULL REFERENCES users(id),
			amount DECIMAL(10, 2) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS expense_splits (
			id SERIAL PRIMARY KEY,
			expense_id INTEGER NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id),
			amount DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS balances (
			id SERIAL PRIMARY KEY,
			group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			amount DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS settlements (
			id SERIAL PRIMARY KEY,
			group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
			from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			amount DECIMAL(10, 2) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_groups_creator_id ON groups(creator_id)`,
		`CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON group_members(group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_expenses_group_id ON expenses(group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_expenses_paid_by_id ON expenses(paid_by_id)`,
		`CREATE INDEX IF NOT EXISTS idx_expense_splits_expense_id ON expense_splits(expense_id)`,
		`CREATE INDEX IF NOT EXISTS idx_expense_splits_user_id ON expense_splits(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_settlements_group_id ON settlements(group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_settlements_from_user_id ON settlements(from_user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_settlements_to_user_id ON settlements(to_user_id)`,
	}

	queries = append(queries, indexQueries...)

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating table: %v", err)
		}
	}

	log.Println("Database tables initialized")
}
