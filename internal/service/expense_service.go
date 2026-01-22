package service

import (
	"fmt"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/repositorypg"
)

type ExpenseService struct {
	userRepo    *repositorypg.UserRepositoryPG
	expenseRepo *repositorypg.ExpenseRepositoryPG
	splitRepo   *repositorypg.ExpenseSplitRepositoryPG
}

func NewExpenseService(
	userRepo *repositorypg.UserRepositoryPG,
	expenseRepo *repositorypg.ExpenseRepositoryPG,
	splitRepo *repositorypg.ExpenseSplitRepositoryPG,
) *ExpenseService {
	return &ExpenseService{
		userRepo:    userRepo,
		expenseRepo: expenseRepo,
		splitRepo:   splitRepo,
	}
}

func (s *ExpenseService) CreateExpense(groupID, paidByID int, amount float64, description string) (*model.ExpenseResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Get user details
	user, err := s.userRepo.GetUserByID(paidByID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	expense := &model.Expense{
		GroupID:     groupID,
		PaidByID:    paidByID,
		Amount:      amount,
		Description: description,
	}

	createdExpense, err := s.expenseRepo.CreateExpense(expense)
	if err != nil {
		return nil, err
	}

	return &model.ExpenseResponse{
		ID:          createdExpense.ID,
		GroupID:     createdExpense.GroupID,
		PaidByID:    createdExpense.PaidByID,
		PaidByName:  user.Name,
		Amount:      createdExpense.Amount,
		Description: createdExpense.Description,
		CreatedAt:   createdExpense.CreatedAt,
	}, nil
}

func (s *ExpenseService) GetExpenseByID(id int) (*model.ExpenseResponse, error) {
	expense, err := s.expenseRepo.GetExpenseByID(id)
	if err != nil {
		return nil, err
	}

	// Get user details
	user, err := s.userRepo.GetUserByID(expense.PaidByID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &model.ExpenseResponse{
		ID:          expense.ID,
		GroupID:     expense.GroupID,
		PaidByID:    expense.PaidByID,
		PaidByName:  user.Name,
		Amount:      expense.Amount,
		Description: expense.Description,
		CreatedAt:   expense.CreatedAt,
	}, nil
}

func (s *ExpenseService) GetExpensesByGroupID(groupID int) ([]*model.ExpenseResponse, error) {
	expenses, err := s.expenseRepo.GetExpensesByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	var responses []*model.ExpenseResponse
	for _, expense := range expenses {
		// Get user details
		user, err := s.userRepo.GetUserByID(expense.PaidByID)
		if err != nil {
			continue // Skip if user not found
		}

		responses = append(responses, &model.ExpenseResponse{
			ID:          expense.ID,
			GroupID:     expense.GroupID,
			PaidByID:    expense.PaidByID,
			PaidByName:  user.Name,
			Amount:      expense.Amount,
			Description: expense.Description,
			CreatedAt:   expense.CreatedAt,
		})
	}

	return responses, nil
}

func (s *ExpenseService) GetExpensesByUserID(userID int) ([]*model.ExpenseResponse, error) {
	expenses, err := s.expenseRepo.GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []*model.ExpenseResponse
	for _, expense := range expenses {
		// Get user details
		user, err := s.userRepo.GetUserByID(expense.PaidByID)
		if err != nil {
			continue // Skip if user not found
		}

		responses = append(responses, &model.ExpenseResponse{
			ID:          expense.ID,
			GroupID:     expense.GroupID,
			PaidByID:    expense.PaidByID,
			PaidByName:  user.Name,
			Amount:      expense.Amount,
			Description: expense.Description,
			CreatedAt:   expense.CreatedAt,
		})
	}

	return responses, nil
}

func (s *ExpenseService) UpdateExpense(id int, amount float64, description string) (*model.ExpenseResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	expense := &model.Expense{
		ID:          id,
		Amount:      amount,
		Description: description,
	}

	updatedExpense, err := s.expenseRepo.UpdateExpense(expense)
	if err != nil {
		return nil, err
	}

	// Get user details
	user, err := s.userRepo.GetUserByID(updatedExpense.PaidByID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &model.ExpenseResponse{
		ID:          updatedExpense.ID,
		GroupID:     updatedExpense.GroupID,
		PaidByID:    updatedExpense.PaidByID,
		PaidByName:  user.Name,
		Amount:      updatedExpense.Amount,
		Description: updatedExpense.Description,
		CreatedAt:   updatedExpense.CreatedAt,
	}, nil
}

func (s *ExpenseService) DeleteExpense(id int) error {
	// Delete all splits for this expense first
	err := s.splitRepo.DeleteSplitsByExpenseID(id)
	if err != nil {
		return err
	}

	return s.expenseRepo.DeleteExpense(id)
}

// Add Split
func (s *ExpenseService) AddSplit(expenseID, userID int, amount float64) (*model.ExpenseSplitResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("split amount must be greater than 0")
	}

	split := &model.ExpenseSplit{
		ExpenseID: expenseID,
		UserID:    userID,
		Amount:    amount,
	}

	createdSplit, err := s.splitRepo.CreateSplit(split)
	if err != nil {
		return nil, err
	}

	return &model.ExpenseSplitResponse{
		ID:        createdSplit.ID,
		ExpenseID: createdSplit.ExpenseID,
		UserID:    createdSplit.UserID,
		Amount:    createdSplit.Amount,
	}, nil
}

func (s *ExpenseService) GetSplitsByExpenseID(expenseID int) ([]*model.ExpenseSplitResponse, error) {
	splits, err := s.splitRepo.GetSplitsByExpenseID(expenseID)
	if err != nil {
		return nil, err
	}

	var responses []*model.ExpenseSplitResponse
	for _, split := range splits {
		responses = append(responses, &model.ExpenseSplitResponse{
			ID:        split.ID,
			ExpenseID: split.ExpenseID,
			UserID:    split.UserID,
			Amount:    split.Amount,
		})
	}

	return responses, nil
}

func (s *ExpenseService) GetSplitsByUserID(userID int) ([]*model.ExpenseSplitResponse, error) {
	splits, err := s.splitRepo.GetSplitsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []*model.ExpenseSplitResponse
	for _, split := range splits {
		responses = append(responses, &model.ExpenseSplitResponse{
			ID:        split.ID,
			ExpenseID: split.ExpenseID,
			UserID:    split.UserID,
			Amount:    split.Amount,
		})
	}

	return responses, nil
}

func (s *ExpenseService) UpdateSplit(id int, amount float64) (*model.ExpenseSplitResponse, error) {
	if amount <= 0 {
		return nil, fmt.Errorf("split amount must be greater than 0")
	}

	split := &model.ExpenseSplit{
		ID:     id,
		Amount: amount,
	}

	updatedSplit, err := s.splitRepo.UpdateSplit(split)
	if err != nil {
		return nil, err
	}

	return &model.ExpenseSplitResponse{
		ID:        updatedSplit.ID,
		ExpenseID: updatedSplit.ExpenseID,
		UserID:    updatedSplit.UserID,
		Amount:    updatedSplit.Amount,
	}, nil
}
