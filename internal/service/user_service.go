package service

import (
	"fmt"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/repositorypg"
)

type UserService struct {
	repo *repositorypg.UserRepositoryPG
}

func NewUserService(repo *repositorypg.UserRepositoryPG) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(email, name string) (*model.UserResponse, error) {
	if email == "" || name == "" {
		return nil, fmt.Errorf("email and name are required")
	}

	// Check if user already exists
	existing, err := s.repo.GetUserByEmail(email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user with email already exists")
	}

	user := &model.User{
		Email: email,
		Name:  name,
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		Name:      createdUser.Name,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (s *UserService) LoginUser(email string) (*model.UserResponse, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) GetUserByID(id int) (*model.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) GetAllUsers() ([]*model.UserResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var responses []*model.UserResponse
	for _, user := range users {
		responses = append(responses, &model.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		})
	}

	return responses, nil
}

func (s *UserService) UpdateUser(id int, name string) (*model.UserResponse, error) {
	user := &model.User{
		ID:   id,
		Name: name,
	}

	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Name:      updatedUser.Name,
		CreatedAt: updatedUser.CreatedAt,
	}, nil
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
