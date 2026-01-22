package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	balances, err := h.balanceService.GetGroupBalancesWithNames(groupID, h.userRepo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balances)
}
