package repositorypg

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
)

type GroupRepositoryPG struct {
	DB *sql.DB
}

func NewGroupRepositoryPG(db *sql.DB) *GroupRepositoryPG {
	return &GroupRepositoryPG{DB: db}
}

func (r *GroupRepositoryPG) CreateGroup(group *model.Group) (*model.Group, error) {
	query := `
		INSERT INTO groups (name, description, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, creator_id, created_at, updated_at
	`

	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		group.Name,
		group.Description,
		group.CreatorID,
		group.CreatedAt,
		group.UpdatedAt,
	).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID, &group.CreatedAt, &group.UpdatedAt)

	if err != nil {
		log.Printf("Error creating group: %v", err)
		return nil, err
	}

	return group, nil
}

func (r *GroupRepositoryPG) GetGroupByID(id int) (*model.Group, error) {
	query := `
		SELECT id, name, description, creator_id, created_at, updated_at
		FROM groups
		WHERE id = $1
	`

	group := &model.Group{}
	err := r.DB.QueryRow(query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&group.CreatorID,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		log.Printf("Error getting group by ID: %v", err)
		return nil, err
	}

	return group, nil
}

func (r *GroupRepositoryPG) GetAllGroups() ([]*model.Group, error) {
	query := `
		SELECT id, name, description, creator_id, created_at, updated_at
		FROM groups
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("Error getting all groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		group := &model.Group{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.CreatorID,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning group: %v", err)
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating groups: %v", err)
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepositoryPG) GetGroupsByUserID(userID int) ([]*model.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.creator_id, g.created_at, g.updated_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
		ORDER BY g.created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error getting groups by user ID: %v", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*model.Group
	for rows.Next() {
		group := &model.Group{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.CreatorID,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning group: %v", err)
			return nil, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating groups: %v", err)
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepositoryPG) UpdateGroup(group *model.Group) (*model.Group, error) {
	query := `
		UPDATE groups
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, description, creator_id, created_at, updated_at
	`

	group.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		group.Name,
		group.Description,
		group.UpdatedAt,
		group.ID,
	).Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID, &group.CreatedAt, &group.UpdatedAt)

	if err != nil {
		log.Printf("Error updating group: %v", err)
		return nil, err
	}

	return group, nil
}

func (r *GroupRepositoryPG) DeleteGroup(id int) error {
	query := `DELETE FROM groups WHERE id = $1`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting group: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}
