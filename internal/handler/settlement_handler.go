package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shreyansh/expense-go-collab-backend/internal/model"
	"github.com/shreyansh/expense-go-collab-backend/internal/service"
)

type SettlementHandler struct {
	settlementService *service.SettlementService
}

func NewSettlementHandler(settlementService *service.SettlementService) *SettlementHandler {
	return &SettlementHandler{
		settlementService: settlementService,
	}
}

// CreateSettlement creates a new payment settlement
// POST /api/settle
func (h *SettlementHandler) CreateSettlement(c *gin.Context) {
	var req model.SettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settlement, err := h.settlementService.CreateSettlement(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, settlement)
}

// GetSettlementByID retrieves a settlement by ID
// GET /api/settle/:id
func (h *SettlementHandler) GetSettlementByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid settlement id"})
		return
	}

	settlement, err := h.settlementService.GetSettlementByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settlement)
}

// GetGroupSettlements retrieves all settlements in a group
// GET /api/settle/group/:group_id
func (h *SettlementHandler) GetGroupSettlements(c *gin.Context) {
	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	settlements, err := h.settlementService.GetSettlementsByGroupID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settlements)
}

// GetUserSettlements retrieves all settlements for a user
// GET /api/settle/user/:user_id
func (h *SettlementHandler) GetUserSettlements(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	settlements, err := h.settlementService.GetSettlementsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settlements": settlements})
}

// GetAllSettlements retrieves all settlements
// GET /api/settle
func (h *SettlementHandler) GetAllSettlements(c *gin.Context) {
	settlements, err := h.settlementService.GetAllSettlements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settlements": settlements})
}
