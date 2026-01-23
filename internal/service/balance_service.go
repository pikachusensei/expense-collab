package service

import (
	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/repositorypg"
)

type BalanceService struct {
	balanceRepo *repositorypg.BalanceRepositoryPG
}

func NewBalanceService(balanceRepo *repositorypg.BalanceRepositoryPG) *BalanceService {
	return &BalanceService{balanceRepo: balanceRepo}
}

func (s *BalanceService) GetUserBalance(userID, groupID int) (float64, error) {
	return s.balanceRepo.GetUserBalanceInGroup(userID, groupID)
}

func (s *BalanceService) GetGroupBalances(groupID int) (map[int]float64, error) {
	return s.balanceRepo.GetGroupBalances(groupID)
}

func (s *BalanceService) GetUserRelativeBalances(
	groupID int,
	currentUserID int,
	userRepo *repositorypg.UserRepositoryPG,
) ([]*model.UserBalanceView, error) {

	balances, err := s.balanceRepo.GetGroupBalances(groupID)
	if err != nil {
		return nil, err
	}

	currentUserBalance := balances[currentUserID]

	var responses []*model.UserBalanceView

	for userID, amount := range balances {

		// 🚫 Rule 1: skip self
		if userID == currentUserID {
			continue
		}

		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			continue
		}

		diff := currentUserBalance - amount

		if diff == 0 {
			continue
		}

		view := &model.UserBalanceView{
			UserID:   userID,
			UserName: user.Name,
			Amount:   abs(diff),
		}

		if diff < 0 {
			view.Type = "you_owe"
		} else {
			view.Type = "owes_you"
		}

		responses = append(responses, view)
	}

	return responses, nil
}
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}


func (s *BalanceService) GetGroupBalancesWithNames(groupID int, userRepo *repositorypg.UserRepositoryPG) ([]*model.UserBalanceResponse, error) {
	balances, err := s.balanceRepo.GetGroupBalances(groupID)
	if err != nil {
		return nil, err
	}

	var responses []*model.UserBalanceResponse
	for userID, amount := range balances {
		user, err := userRepo.GetUserByID(userID)
		if err != nil {
			continue
		}

		responses = append(responses, &model.UserBalanceResponse{
			UserID:   userID,
			UserName: user.Name,
			Amount:   amount,
		})
	}

	return responses, nil
}
