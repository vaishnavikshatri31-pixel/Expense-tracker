package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/expense-tracker/internal/models"
	"github.com/user/expense-tracker/internal/repository"
)

type GroupHandler struct {
	repo *repository.Repository
}

func NewGroupHandler(repo *repository.Repository) *GroupHandler {
	return &GroupHandler{repo: repo}
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if group.Name == "" || group.CreatedBy == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and created_by are required"})
		return
	}

	err := h.repo.CreateGroup(&group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
		return
	}

	c.JSON(http.StatusCreated, group)
}

func (h *GroupHandler) AddMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req struct {
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByPhone(req.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user with this phone number not found"})
		return
	}

	err = h.repo.AddMemberToGroup(groupID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add member to group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member added successfully"})
}
