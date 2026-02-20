package api

import (
	"github.com/gin-gonic/gin"
	"github.com/user/expense-tracker/internal/api/handlers"
	"github.com/user/expense-tracker/internal/repository"
	"github.com/user/expense-tracker/internal/service"
)

func SetupRouter(repo *repository.Repository) *gin.Engine {
	r := gin.Default()

	userHandler := handlers.NewUserHandler(repo)
	groupHandler := handlers.NewGroupHandler(repo)
	expenseService := service.NewExpenseService(repo)
	expenseHandler := handlers.NewExpenseHandler(expenseService)

	// Root / Health Check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Expense Tracker API!",
			"status":  "running",
		})
	})

	// User Routes
	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUser)

	// Group Routes
	r.POST("/groups", groupHandler.CreateGroup)
	r.POST("/groups/:id/members", groupHandler.AddMember)

	// Expense and Settlement Routes
	r.POST("/expenses", expenseHandler.CreateExpense)
	r.GET("/groups/:id/settlements", expenseHandler.GetSettlements)
	r.POST("/groups/:id/settle", expenseHandler.SettleGroup)

	return r
}
