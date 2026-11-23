package routes

import (
	"database/sql"

	"taskmanagement/config"
	"taskmanagement/internal/handler"
	"taskmanagement/internal/middleware"
	"taskmanagement/internal/repository"
	"taskmanagement/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB, authCfg *config.AuthConfig) *gin.Engine {
	router := gin.Default()

	// Health Check
	router.GET("/health", func(c *gin.Context) {
		if err := db.PingContext(c.Request.Context()); err != nil {
			c.JSON(503, gin.H{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Task Management System API is running",
		})
	})

	// Repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	adminDashboardRepo := repository.NewAdminDashboardRepository(db)
	userDashboardRepo := repository.NewUserDashboardRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, authCfg)
	categoryService := service.NewCategoryService(categoryRepo)
	taskService := service.NewTaskService(taskRepo, categoryRepo)
	adminUserService := service.NewAdminUserService(userRepo)
	adminTaskService := service.NewAdminTaskService(taskRepo)
	adminDashboardService := service.NewAdminDashboardService(adminDashboardRepo)
	userDashboardService := service.NewUserDashboardService(userDashboardRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	taskHandler := handler.NewTaskHandler(taskService)
	adminUserHandler := handler.NewAdminUserHandler(adminUserService)
	adminTaskHandler := handler.NewAdminTaskHandler(adminTaskService)
	adminDashboardHandler := handler.NewAdminDashboardHandler(adminDashboardService)
	userDashboardHandler := handler.NewUserDashboardHandler(userDashboardService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(userRepo, authCfg)

	// Public Authentication Routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Protected Authentication Routes
	protectedAuth := router.Group("/auth")
	protectedAuth.Use(authMiddleware.Authenticate())
	{
		protectedAuth.GET("/me", authHandler.Me)
	}

	// Protected User Routes
	userRoutes := router.Group("/api/v1")
	userRoutes.Use(authMiddleware.Authenticate())
	userRoutes.Use(middleware.RequireRole("USER"))
	{
		// User Dashboard
		userRoutes.GET("/dashboard", userDashboardHandler.GetStats)

		// Category CRUD
		categories := userRoutes.Group("/categories")
		{
			categories.POST("", categoryHandler.Create)
			categories.GET("", categoryHandler.FindAll)
			categories.GET("/:id", categoryHandler.FindByID)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		// Task CRUD
		tasks := userRoutes.Group("/tasks")
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.FindAll)
			tasks.GET("/:id", taskHandler.FindByID)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}
	}

	// Protected Admin Routes
	adminRoutes := router.Group("/api/v1/admin")
	adminRoutes.Use(authMiddleware.Authenticate())
	adminRoutes.Use(middleware.RequireRole("ADMIN"))

	adminRoutes.PUT("/users/:id/role", adminUserHandler.UpdateRole)
	adminRoutes.PUT("/users/:id/status", adminUserHandler.UpdateStatus)
	adminRoutes.GET("/dashboard", adminDashboardHandler.GetStats)

	{
		adminRoutes.GET("/tasks", adminTaskHandler.ListTasks)

		users := adminRoutes.Group("/users")
		{
			users.GET("", adminUserHandler.ListUsers)
			users.GET("/:id", adminUserHandler.FindByID)
		}
	}

	return router
}
