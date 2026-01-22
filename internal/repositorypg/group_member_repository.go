package repositorypg

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
)

type GroupMemberRepositoryPG struct {
	DB *sql.DB
}

func NewGroupMemberRepositoryPG(db *sql.DB) *GroupMemberRepositoryPG {
	return &GroupMemberRepositoryPG{DB: db}
}

func (r *GroupMemberRepositoryPG) AddMember(member *model.GroupMember) (*model.GroupMember, error) {
	query := `
		INSERT INTO group_members (group_id, user_id, added_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, group_id, user_id, added_at, updated_at
	`

	member.AddedAt = time.Now()
	member.UpdatedAt = time.Now()

	err := r.DB.QueryRow(
		query,
		member.GroupID,
		member.UserID,
		member.AddedAt,
		member.UpdatedAt,
	).Scan(&member.ID, &member.GroupID, &member.UserID, &member.AddedAt, &member.UpdatedAt)

	if err != nil {
		log.Printf("Error adding member: %v", err)
		return nil, err
	}

	return member, nil
}

func (r *GroupMemberRepositoryPG) RemoveMember(groupID, userID int) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`

	result, err := r.DB.Exec(query, groupID, userID)
	if err != nil {
		log.Printf("Error removing member: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}

func (r *GroupMemberRepositoryPG) GetGroupMembers(groupID int) ([]*model.GroupMember, error) {
	query := `
		SELECT id, group_id, user_id, added_at, updated_at
		FROM group_members
		WHERE group_id = $1
		ORDER BY added_at DESC
	`

	rows, err := r.DB.Query(query, groupID)
	if err != nil {
		log.Printf("Error getting group members: %v", err)
		return nil, err
	}
	defer rows.Close()

	var members []*model.GroupMember
	for rows.Next() {
		member := &model.GroupMember{}
		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&member.AddedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning member: %v", err)
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating members: %v", err)
		return nil, err
	}

	return members, nil
}

func (r *GroupMemberRepositoryPG) GetUserGroups(userID int) ([]*model.GroupMember, error) {
	query := `
		SELECT id, group_id, user_id, added_at, updated_at
		FROM group_members
		WHERE user_id = $1
		ORDER BY added_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error getting user groups: %v", err)
		return nil, err
	}
	defer rows.Close()

	var members []*model.GroupMember
	for rows.Next() {
		member := &model.GroupMember{}
		err := rows.Scan(
			&member.ID,
			&member.GroupID,
			&member.UserID,
			&member.AddedAt,
			&member.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning member: %v", err)
			return nil, err
		}
		members = append(members, member)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating members: %v", err)
		return nil, err
	}

	return members, nil
}

func (r *GroupMemberRepositoryPG) IsMember(groupID, userID int) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = $1 AND user_id = $2
	`

	var count int
	err := r.DB.QueryRow(query, groupID, userID).Scan(&count)
	if err != nil {
		log.Printf("Error checking membership: %v", err)
		return false, err
	}

	return count > 0, nil
}
