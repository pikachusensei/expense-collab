package service

import (
	"fmt"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/repositorypg"
)

type GroupService struct {
	userRepo    *repositorypg.UserRepositoryPG
	groupRepo   *repositorypg.GroupRepositoryPG
	memberRepo  *repositorypg.GroupMemberRepositoryPG
	expenseRepo *repositorypg.ExpenseRepositoryPG
	splitRepo   *repositorypg.ExpenseSplitRepositoryPG
	balanceRepo *repositorypg.BalanceRepositoryPG
}

func NewGroupService(
	userRepo *repositorypg.UserRepositoryPG,
	groupRepo *repositorypg.GroupRepositoryPG,
	memberRepo *repositorypg.GroupMemberRepositoryPG,
	expenseRepo *repositorypg.ExpenseRepositoryPG,
	splitRepo *repositorypg.ExpenseSplitRepositoryPG,
	balanceRepo *repositorypg.BalanceRepositoryPG,
) *GroupService {
	return &GroupService{
		userRepo:    userRepo,
		groupRepo:   groupRepo,
		memberRepo:  memberRepo,
		expenseRepo: expenseRepo,
		splitRepo:   splitRepo,
		balanceRepo: balanceRepo,
	}
}

func (s *GroupService) CreateGroup(name string, description string, creatorID int) (*model.GroupResponse, error) {
	if name == "" {
		return nil, fmt.Errorf("group name is required")
	}

	group := &model.Group{
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
	}

	createdGroup, err := s.groupRepo.CreateGroup(group)
	if err != nil {
		return nil, err
	}

	// Add creator as member
	member := &model.GroupMember{
		GroupID: createdGroup.ID,
		UserID:  creatorID,
	}
	s.memberRepo.AddMember(member)

	return &model.GroupResponse{
		ID:          createdGroup.ID,
		Name:        createdGroup.Name,
		Description: createdGroup.Description,
		CreatorID:   createdGroup.CreatorID,
		CreatedAt:   createdGroup.CreatedAt,
	}, nil
}

func (s *GroupService) GetGroupByID(id int) (*model.GroupResponse, error) {
	group, err := s.groupRepo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}

	return &model.GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		CreatorID:   group.CreatorID,
		CreatedAt:   group.CreatedAt,
	}, nil
}

func (s *GroupService) GetAllGroups() ([]*model.GroupResponse, error) {
	groups, err := s.groupRepo.GetAllGroups()
	if err != nil {
		return nil, err
	}

	var responses []*model.GroupResponse
	for _, group := range groups {
		responses = append(responses, &model.GroupResponse{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			CreatorID:   group.CreatorID,
			CreatedAt:   group.CreatedAt,
		})
	}

	return responses, nil
}

func (s *GroupService) GetGroupsByUserID(userID int) ([]*model.GroupResponse, error) {
	groups, err := s.groupRepo.GetGroupsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []*model.GroupResponse
	for _, group := range groups {
		responses = append(responses, &model.GroupResponse{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			CreatorID:   group.CreatorID,
			CreatedAt:   group.CreatedAt,
		})
	}

	return responses, nil
}

func (s *GroupService) UpdateGroup(id int, name, description string) (*model.GroupResponse, error) {
	group := &model.Group{
		ID:          id,
		Name:        name,
		Description: description,
	}

	updatedGroup, err := s.groupRepo.UpdateGroup(group)
	if err != nil {
		return nil, err
	}

	return &model.GroupResponse{
		ID:          updatedGroup.ID,
		Name:        updatedGroup.Name,
		Description: updatedGroup.Description,
		CreatorID:   updatedGroup.CreatorID,
		CreatedAt:   updatedGroup.CreatedAt,
	}, nil
}

func (s *GroupService) DeleteGroup(id int) error {
	return s.groupRepo.DeleteGroup(id)
}

func (s *GroupService) AddMemberToGroup(groupID, userID int) (*model.GroupMemberResponse, error) {
	isMember, err := s.memberRepo.IsMember(groupID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, fmt.Errorf("user is already a member of this group")
	}

	member := &model.GroupMember{
		GroupID: groupID,
		UserID:  userID,
	}

	createdMember, err := s.memberRepo.AddMember(member)
	if err != nil {
		return nil, err
	}

	return &model.GroupMemberResponse{
		ID:      createdMember.ID,
		GroupID: createdMember.GroupID,
		UserID:  createdMember.UserID,
		AddedAt: createdMember.AddedAt,
	}, nil
}

// AddMemberToGroupByEmail adds a member by email and returns enriched response with user details
func (s *GroupService) AddMemberToGroupByEmail(groupID int, email string) (*model.GroupMemberResponse, error) {
	// Look up user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found with email: %s", email)
	}

	// Check if already a member
	isMember, err := s.memberRepo.IsMember(groupID, user.ID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, fmt.Errorf("user is already a member of this group")
	}

	// Create member
	member := &model.GroupMember{
		GroupID: groupID,
		UserID:  user.ID,
	}

	createdMember, err := s.memberRepo.AddMember(member)
	if err != nil {
		return nil, err
	}

	// Return enriched response with user details
	return &model.GroupMemberResponse{
		ID:       createdMember.ID,
		GroupID:  createdMember.GroupID,
		UserID:   createdMember.UserID,
		UserName: user.Name,
		Email:    user.Email,
		AddedAt:  createdMember.AddedAt,
	}, nil
}

func (s *GroupService) RemoveMemberFromGroup(groupID, userID int) error {
	return s.memberRepo.RemoveMember(groupID, userID)
}

func (s *GroupService) GetGroupMembers(groupID int) ([]*model.GroupMemberResponse, error) {
	members, err := s.memberRepo.GetGroupMembersWithDetails(groupID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
