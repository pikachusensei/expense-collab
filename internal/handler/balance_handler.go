package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/repositorypg"
	"github.com/shreyansh/expense-go-collab-backend/internal/service"
)

type BalanceHandler struct {
	balanceService *service.BalanceService
	userRepo       *repositorypg.UserRepositoryPG
}

func NewBalanceHandler(balanceService *service.BalanceService, userRepo *repositorypg.UserRepositoryPG) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
		userRepo:       userRepo,
	}
}

func (h *BalanceHandler) GetUserBalance(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	balance, err := h.balanceService.GetUserBalance(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "group_id": groupID, "balance": balance})
}

func (h *BalanceHandler) GetGroupBalances(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	// Get current user ID from query parameter or header
	userID := c.DefaultQuery("user_id", "0")
	currentUserID, err := strconv.Atoi(userID)
	if err != nil {
		currentUserID = 0
	}

	balances, err := h.balanceService.GetGroupBalancesWithNames(groupID, h.userRepo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter out zero/near-zero balances and the current user's own balance
	var filteredBalances []*model.UserBalanceResponse
	for _, balance := range balances {
		// Skip if balance is near zero or if it's the current user
		if balance.Amount > 0.01 || balance.Amount < -0.01 {
			if currentUserID > 0 && balance.UserID == currentUserID {
				continue // Skip current user's balance relative to themselves
			}
			filteredBalances = append(filteredBalances, balance)
		}
	}

	c.JSON(http.StatusOK, filteredBalances)
}
