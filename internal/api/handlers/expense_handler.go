package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/user/expense-tracker/internal/models"
	"github.com/user/expense-tracker/internal/service"
)

type ExpenseHandler struct {
	service *service.ExpenseService
}

func NewExpenseHandler(service *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

type CreateExpenseRequest struct {
	GroupID     uuid.UUID       `json:"group_id"`
	PaidBy      uuid.UUID       `json:"paid_by"`
	Description string          `json:"description"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	SplitType   string          `json:"split_type"`
	Splits      []struct {
		UserID uuid.UUID       `json:"user_id"`
		Amount decimal.Decimal `json:"amount"`
	} `json:"splits"`
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense := &models.Expense{
		GroupID:     req.GroupID,
		PaidBy:      req.PaidBy,
		Description: req.Description,
		TotalAmount: req.TotalAmount,
		SplitType:   req.SplitType,
	}

	err := h.service.CreateExpense(expense, req.Splits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) GetSettlements(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	transactions, err := h.service.GetSettlements(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *ExpenseHandler) SettleGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	err = h.service.SettleGroup(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group settled and balance reset to zero"})
}
