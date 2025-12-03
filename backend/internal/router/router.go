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
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())

	// ------------------------------------------------------
	// MUST BE BEFORE REGISTERING ANY ROUTES
	// ------------------------------------------------------
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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
	analyticsHandler := handlers.NewAnalyticsHandler(reportService)

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
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
		users.GET("/me", userHandler.GetCurrentUser)
		users.GET("/:id", userHandler.GetUserByID)
		users.PATCH("/:id/password", userHandler.ChangePassword)
		users.PATCH("/:id/project", userHandler.SetDefaultProject)
		users.PATCH("/:id/image", userHandler.ChangeImage)

		admin := users.Group("")
		admin.Use(middleware.RequiredRole("superadmin"))
		{
			admin.GET("", userHandler.ListUsers)
			admin.PATCH("/:id", userHandler.UpdateUser)
			admin.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// ---------- Reports ----------

	reports := protected.Group("/reports")
	{
		reports.GET("/:projectID/csp", reportHandler.ListReportsByProjectID)
	}
	api.POST("/reports/:projectID/csp", reportHandler.CreateReport)

	// ---------- Analytics ----------
	analytics := protected.Group("/analytics")
	{
		analytics.GET("/:projectID/graph-data", analyticsHandler.GetReportGraphData)
	}

	// ---------- Signout ----------
	protected.POST("/signout", authHandler.SignOut)

	return router
}
