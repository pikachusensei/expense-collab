package service

import (
	"fmt"

	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/repository"
)

type SettlementService struct {
	settlementRepo repository.SettlementRepository
	userRepo       repository.UserRepository
	balanceRepo    repository.BalanceRepository
}

func NewSettlementService(
	settlementRepo repository.SettlementRepository,
	userRepo repository.UserRepository,
	balanceRepo repository.BalanceRepository,
) *SettlementService {
	return &SettlementService{
		settlementRepo: settlementRepo,
		userRepo:       userRepo,
		balanceRepo:    balanceRepo,
	}
}

// CreateSettlement creates a new payment settlement and updates user balances
func (s *SettlementService) CreateSettlement(req *model.SettlementRequest) (*model.SettlementResponse, error) {
	// Validate users exist
	fromUser, err := s.userRepo.GetUserByID(req.FromUserID)
	if err != nil {
		return nil, fmt.Errorf("from user not found: %v", err)
	}

	toUser, err := s.userRepo.GetUserByID(req.ToUserID)
	if err != nil {
		return nil, fmt.Errorf("to user not found: %v", err)
	}

	// Validate amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Create settlement record
	settlement := &model.Settlement{
		GroupID:     req.GroupID,
		FromUserID:  req.FromUserID,
		ToUserID:    req.ToUserID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	created, err := s.settlementRepo.CreateSettlement(settlement)
	if err != nil {
		return nil, err
	}

	// Prepare response with user names
	response := &model.SettlementResponse{
		ID:           created.ID,
		GroupID:      created.GroupID,
		FromUserID:   created.FromUserID,
		FromUserName: fromUser.Name,
		ToUserID:     created.ToUserID,
		ToUserName:   toUser.Name,
		Amount:       created.Amount,
		Description:  created.Description,
		CreatedAt:    created.CreatedAt,
	}

	return response, nil
}

// GetSettlementByID retrieves a settlement by ID
func (s *SettlementService) GetSettlementByID(id int) (*model.SettlementResponse, error) {
	settlement, err := s.settlementRepo.GetSettlementByID(id)
	if err != nil {
		return nil, err
	}

	fromUser, _ := s.userRepo.GetUserByID(settlement.FromUserID)
	toUser, _ := s.userRepo.GetUserByID(settlement.ToUserID)

	response := &model.SettlementResponse{
		ID:           settlement.ID,
		GroupID:      settlement.GroupID,
		FromUserID:   settlement.FromUserID,
		FromUserName: fromUser.Name,
		ToUserID:     settlement.ToUserID,
		ToUserName:   toUser.Name,
		Amount:       settlement.Amount,
		Description:  settlement.Description,
		CreatedAt:    settlement.CreatedAt,
	}

	return response, nil
}

// GetSettlementsByGroupID retrieves all settlements in a group
func (s *SettlementService) GetSettlementsByGroupID(groupID int) (*model.GroupSettlementResponse, error) {
	settlements, err := s.settlementRepo.GetSettlementsByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	var responses []model.SettlementResponse
	totalAmount := 0.0

	for _, settlement := range settlements {
		fromUser, _ := s.userRepo.GetUserByID(settlement.FromUserID)
		toUser, _ := s.userRepo.GetUserByID(settlement.ToUserID)

		response := model.SettlementResponse{
			ID:           settlement.ID,
			GroupID:      settlement.GroupID,
			FromUserID:   settlement.FromUserID,
			FromUserName: fromUser.Name,
			ToUserID:     settlement.ToUserID,
			ToUserName:   toUser.Name,
			Amount:       settlement.Amount,
			Description:  settlement.Description,
			CreatedAt:    settlement.CreatedAt,
		}
		responses = append(responses, response)
		totalAmount += settlement.Amount
	}

	return &model.GroupSettlementResponse{
		Settlements: responses,
		TotalAmount: totalAmount,
	}, nil
}

// GetSettlementsByUserID retrieves all settlements for a user
func (s *SettlementService) GetSettlementsByUserID(userID int) ([]*model.SettlementResponse, error) {
	settlements, err := s.settlementRepo.GetSettlementsByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []*model.SettlementResponse

	for _, settlement := range settlements {
		fromUser, _ := s.userRepo.GetUserByID(settlement.FromUserID)
		toUser, _ := s.userRepo.GetUserByID(settlement.ToUserID)

		response := &model.SettlementResponse{
			ID:           settlement.ID,
			GroupID:      settlement.GroupID,
			FromUserID:   settlement.FromUserID,
			FromUserName: fromUser.Name,
			ToUserID:     settlement.ToUserID,
			ToUserName:   toUser.Name,
			Amount:       settlement.Amount,
			Description:  settlement.Description,
			CreatedAt:    settlement.CreatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// GetAllSettlements retrieves all settlements
func (s *SettlementService) GetAllSettlements() ([]*model.SettlementResponse, error) {
	settlements, err := s.settlementRepo.GetAllSettlements()
	if err != nil {
		return nil, err
	}

	var responses []*model.SettlementResponse

	for _, settlement := range settlements {
		fromUser, _ := s.userRepo.GetUserByID(settlement.FromUserID)
		toUser, _ := s.userRepo.GetUserByID(settlement.ToUserID)

		response := &model.SettlementResponse{
			ID:           settlement.ID,
			GroupID:      settlement.GroupID,
			FromUserID:   settlement.FromUserID,
			FromUserName: fromUser.Name,
			ToUserID:     settlement.ToUserID,
			ToUserName:   toUser.Name,
			Amount:       settlement.Amount,
			Description:  settlement.Description,
			CreatedAt:    settlement.CreatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
