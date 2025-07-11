package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"web-clean/domain"
	"web-clean/infra"
	"web-clean/infra/database"
	byjson "web-clean/infra/loader/json"
	"web-clean/infra/web"
	oldRepository "web-clean/repository"

	// Clean Architecture layers
	"web-clean/internal/application/service"
	"web-clean/internal/infrastructure/repository"
	userHttpHandler "web-clean/internal/interface/http"
)

func main() {
	// Initialize infrastructure context
	context, err := infra.Prepare(infra.PrepareConfig{Loader: byjson.JSONLoader})
	if err != nil {
		panic(err)
	}

	// Initialize database
	db, err := database.From(context)
	if err != nil {
		panic(err)
	}

	// Auto-migrate schemas (including new user schema)
	err = database.AutoMigrateRegisteredSchema(db)
	if err != nil {
		panic(err)
	}

	// Initialize Clean Architecture layers following dependency inversion principle

	// Infrastructure Layer - implements domain interfaces
	userRepo := repository.NewUserRepository(db)

	// Application Layer - contains business logic
	userService := service.NewUserService(userRepo, context.Log)

	// Interface Layer - handles HTTP concerns
	userHandler := userHttpHandler.NewUserHandler(userService, context.Log)

	// Legacy components (keeping for existing functionality)
	logsPersister := oldRepository.Logs{
		Context:  context,
		Database: db,
	}

	contextMiddleware := web.ContextMiddleware(func(log domain.Log) *web.Context {
		return &web.Context{
			Database: db,
			Log:      log,
		}
	}, context.Log, &logsPersister)

	errorsPersister := oldRepository.Errors{
		Context:          context,
		FallbackFilePath: "./errors",
		Database:         db,
	}

	// Initialize web server with Clean Architecture routes
	server := web.Gin(context, func(engine *gin.Engine) {
		// Global middleware
		engine.Use(web.RequestIDMiddleware(func() string {
			return uuid.NewString()
		}))

		engine.Use(web.ErrorPersisterMiddleware(errorsPersister, context.Log, web.RequestIdGetter))

		engine.Use(web.Recover(func(context *gin.Context, err any) {
			// Handle panics gracefully
			context.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal_server_error",
				"message": "An internal error occurred",
			})
		}))

		engine.Use(contextMiddleware)

		// Health check endpoint
		engine.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "web-clean",
			})
		})

		// API v1 routes following Clean Architecture
		apiV1 := engine.Group("/api/v1")
		{
			// User management endpoints
			users := apiV1.Group("/users")
			{
				users.POST("", userHandler.CreateUser)           // POST /api/v1/users
				users.GET("", userHandler.ListUsers)             // GET /api/v1/users?offset=0&limit=10
				users.GET("/:id", userHandler.GetUserByID)       // GET /api/v1/users/:id
				users.PUT("/:id", userHandler.UpdateUserProfile) // PUT /api/v1/users/:id
				users.DELETE("/:id", userHandler.DeleteUser)     // DELETE /api/v1/users/:id
			}
		}

		// API documentation endpoint
		apiV1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Clean Architecture API v1",
				"endpoints": gin.H{
					"users": gin.H{
						"POST /api/v1/users":       "Create a new user",
						"GET /api/v1/users":        "List users with pagination",
						"GET /api/v1/users/:id":    "Get user by ID",
						"PUT /api/v1/users/:id":    "Update user profile",
						"DELETE /api/v1/users/:id": "Delete user",
					},
					"health": "GET /health - Health check",
				},
			})
		})
	})

	// Start the server
	context.Log.Infow("Starting Clean Architecture web server",
		"architecture", "Clean Architecture",
		"layers", []string{"Domain", "Application", "Infrastructure", "Interface"},
		"patterns", []string{"Dependency Inversion", "Separation of Concerns", "Single Responsibility"},
	)

	server.Serve()
}
