package app

import (
	"math/rand"
	"net/http"
	"time"

	"core-gateway/internal/domain/auth"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// legacy types for backwards compatibility
type systemStats struct {
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Storage     int    `json:"storage"`
	Temperature int    `json:"temperature"`
	Uptime      string `json:"uptime"`
	Network     struct {
		Up   string `json:"up"`
		Down string `json:"down"`
	} `json:"network"`
}

type serviceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// RegisterRoutes sets up all API routes
func (a *Application) RegisterRoutes() {
	// swagger documentation
	a.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// health endpoint
	a.Echo.GET("/health", a.healthCheck)

	// legacy endpoints (for backwards compatibility)
	a.Echo.GET("/api/system/stats", a.getSystemStats)
	a.Echo.GET("/api/services", a.getLegacyServices)

	// API v1
	api := a.Echo.Group("/api")
	{
		// auth routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", a.AuthHandler.Register)
			authGroup.POST("/login", a.AuthHandler.Login)
			authGroup.POST("/logout", a.AuthHandler.Logout)
			authGroup.POST("/demo", a.AuthHandler.LoginAsDemo)
			authGroup.GET("/me", a.AuthHandler.Me)
		}

		// user management routes (admin only, except get/update own)
		users := api.Group("/users")
		users.Use(auth.RequireAuth())
		{
			users.GET("", a.UserHandler.ListUsers, auth.RequireRole(auth.RoleAdmin))
			users.POST("", a.UserHandler.CreateUser, auth.RequireRole(auth.RoleAdmin))
			users.GET("/:id", a.UserHandler.GetUser)    // auth check in handler (admin or self)
			users.PUT("/:id", a.UserHandler.UpdateUser) // auth check in handler (admin or self)
			users.DELETE("/:id", a.UserHandler.DeleteUser, auth.RequireRole(auth.RoleAdmin))
		}

		// service monitoring routes
		services := api.Group("/services")
		{
			services.POST("", a.ServiceHandler.CreateService)
			services.GET("/list", a.ServiceHandler.ListServices)
			services.GET("/:id", a.ServiceHandler.GetService)
			services.PUT("/:id", a.ServiceHandler.UpdateService)
			services.DELETE("/:id", a.ServiceHandler.DeleteService)
			services.GET("/:id/history", a.ServiceHandler.GetServiceHistory)
			services.GET("/:id/stats", a.ServiceHandler.GetServiceStats)
			services.GET("/stats/all", a.ServiceHandler.GetAllServicesStats)
		}

		// budget tracking routes
		budget := api.Group("/budget")
		{
			// priority groups (reference data)
			budget.GET("/priority-groups", a.BudgetHandler.GetPriorityGroups)

			// categories
			categories := budget.Group("/categories")
			categories.POST("", a.BudgetHandler.CreateCategory)
			categories.GET("", a.BudgetHandler.ListCategories)
			categories.PUT("/:id", a.BudgetHandler.UpdateCategory)
			categories.DELETE("/:id", a.BudgetHandler.DeleteCategory)

			// tags
			tags := budget.Group("/tags")
			tags.POST("", a.BudgetHandler.CreateTag)
			tags.GET("", a.BudgetHandler.ListTags)
			tags.PUT("/:id", a.BudgetHandler.UpdateTag)
			tags.DELETE("/:id", a.BudgetHandler.DeleteTag)

			// expenses
			expenses := budget.Group("/expenses")
			expenses.POST("", a.BudgetHandler.CreateExpense)
			expenses.GET("", a.BudgetHandler.ListExpenses)
			expenses.GET("/search", a.BudgetHandler.SearchExpenses)
			expenses.GET("/check-skipped", a.BudgetHandler.CheckSkippedExpense)
			expenses.GET("/:id", a.BudgetHandler.GetExpense)
			expenses.PUT("/:id", a.BudgetHandler.UpdateExpense)
			expenses.DELETE("/:id", a.BudgetHandler.DeleteExpense)

			// incomes
			incomes := budget.Group("/incomes")
			incomes.POST("", a.BudgetHandler.CreateIncome)
			incomes.GET("", a.BudgetHandler.ListIncomes)
			incomes.GET("/check-skipped", a.BudgetHandler.CheckSkippedIncome)
			incomes.GET("/occurrences", a.BudgetHandler.GetIncomeOccurrences)
			incomes.GET("/:id", a.BudgetHandler.GetIncome)
			incomes.PUT("/:id", a.BudgetHandler.UpdateIncome)
			incomes.DELETE("/:id", a.BudgetHandler.DeleteIncome)

			// budget remaining
			budget.GET("/remaining", a.BudgetHandler.GetBudgetRemaining)

			// stats
			stats := budget.Group("/stats")
			stats.GET("/summary", a.BudgetHandler.GetSummary)
			stats.GET("/trends", a.BudgetHandler.GetTrends)
			stats.GET("/category-breakdown", a.BudgetHandler.GetCategoryBreakdown)
			stats.GET("/savings-rate", a.BudgetHandler.GetSavingsRate)
			stats.GET("/health-score", a.BudgetHandler.GetHealthScore)

			// New budget analytics endpoints
			budget.GET("/velocity", a.BudgetHandler.GetSpendingVelocity)
			budget.GET("/forecast/upcoming", a.BudgetHandler.GetUpcomingBills)
			budget.GET("/current-total-money", a.BudgetHandler.GetCurrentTotalMoney)

			// Financial Health Analysis
			budget.GET("/analysis/503020", a.BudgetHandler.GetFiftyThirtyTwenty)
			budget.GET("/analysis/weekday-pattern", a.BudgetHandler.GetWeekdayPattern)
			budget.GET("/analysis/merchants", a.BudgetHandler.GetMerchantAnalysis)
			budget.GET("/trends/month-over-month", a.BudgetHandler.GetMonthOverMonthTrends)

			// Subscriptions
			budget.GET("/subscriptions", a.BudgetHandler.GetSubscriptions)

			// Recurring Incomes
			budget.GET("/recurring-incomes", a.BudgetHandler.GetRecurringIncomes)

			// Category budgets
			categoryBudgets := budget.Group("/category-budgets")
			categoryBudgets.POST("", a.BudgetHandler.CreateCategoryBudget)
			categoryBudgets.GET("", a.BudgetHandler.GetCategoryBudgets)
			categoryBudgets.PUT("/:id", a.BudgetHandler.UpdateCategoryBudget)
			categoryBudgets.DELETE("/:id", a.BudgetHandler.DeleteCategoryBudget)
		}
	}
}

// healthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (a *Application) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// getSystemStats godoc
// @Summary Get system statistics
// @Description Get mock system statistics (CPU, memory, etc.)
// @Tags system
// @Produce json
// @Success 200 {object} systemStats
// @Router /api/system/stats [get]
func (a *Application) getSystemStats(c echo.Context) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	stats := systemStats{
		CPU:         rand.Intn(30) + 10,
		Memory:      rand.Intn(40) + 20,
		Storage:     68,
		Temperature: 45,
		Uptime:      "42d 13h 27m",
		Network: struct {
			Up   string `json:"up"`
			Down string `json:"down"`
		}{
			Up:   "125.4 Mbps",
			Down: "342.8 Mbps",
		},
	}

	return c.JSON(http.StatusOK, stats)
}

// getLegacyServices godoc
// @Summary Get legacy services (deprecated)
// @Description Get mock service list for backwards compatibility
// @Tags legacy
// @Produce json
// @Success 200 {array} serviceStatus
// @Deprecated true
// @Router /api/services [get]
func (a *Application) getLegacyServices(c echo.Context) error {
	services := []serviceStatus{
		{Name: "Docker Manager", Type: "Container", Status: "online"},
		{Name: "PostgreSQL", Type: "Database", Status: "online"},
		{Name: "Nextcloud", Type: "Storage", Status: "online"},
		{Name: "Vault", Type: "Security", Status: "online"},
		{Name: "Jellyfin", Type: "Media", Status: "offline"},
		{Name: "Mail Server", Type: "Email", Status: "maintenance"},
	}

	return c.JSON(http.StatusOK, services)
}
