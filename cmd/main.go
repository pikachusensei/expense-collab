package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shreyansh/expense-go-collab-backend/internal/config"
	"github.com/shreyansh/expense-go-collab-backend/internal/handler"
	"github.com/shreyansh/expense-go-collab-backend/internal/repositorypg"
	"github.com/shreyansh/expense-go-collab-backend/internal/service"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize database
	db := config.InitDB()
	defer db.Close()

	// Initialize repositories
	userRepo := repositorypg.NewUserRepositoryPG(db)
	groupRepo := repositorypg.NewGroupRepositoryPG(db)
	memberRepo := repositorypg.NewGroupMemberRepositoryPG(db)
	expenseRepo := repositorypg.NewExpenseRepositoryPG(db)
	splitRepo := repositorypg.NewExpenseSplitRepositoryPG(db)
	balanceRepo := repositorypg.NewBalanceRepositoryPG(db)
	settlementRepo := repositorypg.NewSettlementRepositoryPG(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	groupService := service.NewGroupService(userRepo, groupRepo, memberRepo, expenseRepo, splitRepo, balanceRepo)
	expenseService := service.NewExpenseService(expenseRepo, splitRepo)
	balanceService := service.NewBalanceService(balanceRepo)
	settlementService := service.NewSettlementService(settlementRepo, userRepo, balanceRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	groupHandler := handler.NewGroupHandler(groupService)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	balanceHandler := handler.NewBalanceHandler(balanceService, userRepo)
	settlementHandler := handler.NewSettlementHandler(settlementService)

	// Create router
	router := gin.Default()

	// Metrics middleware
	router.Use(handler.MetricsMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapF(promhttp.Handler().ServeHTTP))

	// User routes
	router.POST("/api/users/register", userHandler.Register)
	router.GET("/api/users/login", userHandler.Login)
	router.GET("/api/users", userHandler.GetAllUsers)
	router.GET("/api/users/:id", userHandler.GetUser)
	router.PUT("/api/users/:id", userHandler.UpdateUser)
	router.DELETE("/api/users/:id", userHandler.DeleteUser)

	// Group routes
	router.POST("/api/groups", groupHandler.CreateGroup)
	router.GET("/api/groups", groupHandler.GetAllGroups) //////////TO REMOVE THIS LATER
	router.GET("/api/groups/:id", groupHandler.GetGroup)
	router.PUT("/api/groups/:id", groupHandler.UpdateGroup)
	router.DELETE("/api/groups/:id", groupHandler.DeleteGroup)
	router.GET("/api/groups/user/:user_id", groupHandler.GetUserGroups)

	// Group member routes (use different path structure)
	router.POST("/api/members", groupHandler.AddGroupMember)
	router.DELETE("/api/members/:group_id/:user_id", groupHandler.RemoveGroupMember)
	router.GET("/api/members/group/:group_id", groupHandler.GetGroupMembers)

	// Expense routes
	router.POST("/api/expenses", expenseHandler.CreateExpense)
	router.GET("/api/expenses", expenseHandler.GetExpense)
	router.GET("/api/expenses/group/:group_id", expenseHandler.GetGroupExpenses)
	router.GET("/api/expenses/user/:user_id", expenseHandler.GetUserExpenses)
	router.PUT("/api/expenses/:id", expenseHandler.UpdateExpense)
	router.DELETE("/api/expenses/:id", expenseHandler.DeleteExpense)

	// Expense split routes
	router.POST("/api/splits", expenseHandler.AddExpenseSplit)
	router.GET("/api/splits/expense/:expense_id", expenseHandler.GetExpenseSplits)
	router.GET("/api/splits/user/:user_id", expenseHandler.GetUserSplits)
	router.PUT("/api/splits/:id", expenseHandler.UpdateExpenseSplit)

	// Balance routes
	router.GET("/api/balance/user/:user_id/group/:group_id", balanceHandler.GetUserBalance)
	router.GET("/api/balance/group/:group_id", balanceHandler.GetGroupBalances)

	// Settlement/Payment routes
	router.POST("/api/settle", settlementHandler.CreateSettlement)
	router.GET("/api/settle/:id", settlementHandler.GetSettlementByID)
	router.GET("/api/settle/group/:group_id", settlementHandler.GetGroupSettlements)
	router.GET("/api/settle/user/:user_id", settlementHandler.GetUserSettlements)
	router.GET("/api/settle", settlementHandler.GetAllSettlements)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
