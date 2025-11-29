package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/config"
	docs "github.com/notFil/cspotlight/docs"
	"github.com/notFil/cspotlight/internal/handlers"
	"github.com/notFil/cspotlight/internal/middleware"
	"github.com/notFil/cspotlight/internal/repositories"
	"github.com/notFil/cspotlight/internal/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetUpRouter(db *gorm.DB, jwtConfig config.JWTConfig) *gin.Engine {
	router := gin.Default()

	// ------------------------------------------------------
	// MUST BE BEFORE REGISTERING ANY ROUTES
	// ------------------------------------------------------
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ------------------------------------------------------
	// Setup Repositories + Services + Handlers
	// ------------------------------------------------------
	projectRepo := repositories.NewProjectRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	userRepo := repositories.NewUserRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	projectService := services.NewProjectService(projectRepo)
	teamService := services.NewTeamService(teamRepo)
	userService := services.NewUserService(userRepo, projectRepo)
	reportService := services.NewReportService(reportRepo, projectRepo)

	authHandler := handlers.NewAuthHandler(userService, jwtConfig)
	projectHandler := handlers.NewProjectHandler(projectService)
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	reportHandler := handlers.NewReportHandler(reportService)

	// ------------------------------------------------------

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"

	// ------------------------------------------------------
	// API Routes
	// ------------------------------------------------------
	api := router.Group("/api/v1")

	// -------------------------
	// Public
	// -------------------------
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// -------------------------
	// Protected (JWT required)
	// -------------------------
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(jwtConfig.SecretKey))

	// ---------- Projects ----------
	projects := protected.Group("/projects")
	{
		projects.GET("", projectHandler.ListProjects) // no trailing slash
		projects.GET("/:id", projectHandler.GetProjectByID)

		admin := projects.Group("")
		admin.Use(middleware.RequiredRole("admin", "superadmin"))
		{
			admin.POST("", projectHandler.CreateProject)
			admin.PUT("/:id", projectHandler.UpdateProject)
			admin.DELETE("/:id", projectHandler.DeleteProject)
		}
	}

	// ---------- Teams ----------
	teams := protected.Group("/teams")
	teams.Use(middleware.RequiredRole("superadmin"))
	{
		teams.GET("", teamHandler.ListTeams)
		teams.GET("/:id", teamHandler.GetTeamByID)
		teams.POST("", teamHandler.CreateTeam)
		teams.PUT("/:id", teamHandler.UpdateTeam)
		teams.DELETE("/:id", teamHandler.DeleteTeam)
	}

	// ---------- Users ----------
	users := protected.Group("/users")
	{
		users.GET("/:id", userHandler.GetUserByID)

		admin := users.Group("")
		admin.Use(middleware.RequiredRole("admin", "superadmin"))
		{
			admin.GET("", userHandler.ListUsers)
			admin.POST("", userHandler.CreateUser)
			admin.PUT("/:id", userHandler.UpdateUser)
			admin.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// ---------- Reports ----------

	reports := protected.Group("/reports")
	{
		reports.GET("/:projectID/csp", reportHandler.ListReportsByProjectID)
	}
	api.POST("/reports/:projectID/csp", reportHandler.CreateReport)

	// ---------- Signout ----------
	protected.POST("/signout", authHandler.SignOut)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
