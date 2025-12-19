package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
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

func SetUpRouter(db *gorm.DB, baseURL string, staticPath string, sessionCfg config.Session, corsCfg config.CORS, redisCfg config.Redis, isProduction bool) *gin.Engine {
	var store sessions.Store
	var err error

	router := gin.New()
	if sessionCfg.UseCookieStore {
		store = cookie.NewStore(sessionCfg.SecretKey)
	} else {
		store, err = redis.NewStore(redisCfg.IdleConns, "tcp", redisCfg.Addr, redisCfg.Username, redisCfg.Password, sessionCfg.SecretKey)
		if err != nil {
			panic("Failed to create redis store: " + err.Error())
		}
	}
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   sessionCfg.ExpiryInSeconds,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
	})
	router.Use(sessions.Sessions("session", store))
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RateLimiter())
	router.Static("/static", staticPath)

	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsCfg.Origins,
		AllowMethods:     corsCfg.AllowMethods,
		AllowHeaders:     corsCfg.AllowHeaders,
		ExposeHeaders:    corsCfg.ExposeHeaders,
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	// ------------------------------------------------------
	// Setup Repositories + Services + Handlers
	// ------------------------------------------------------
	projectRepo := repositories.NewProjectRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	userRepo := repositories.NewUserRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	projectService := services.NewProjectService(projectRepo, baseURL)
	teamService := services.NewTeamService(teamRepo)
	userService := services.NewUserService(userRepo)
	reportService := services.NewReportService(reportRepo)

	authHandler := handlers.NewAuthHandler(userService)
	projectHandler := handlers.NewProjectHandler(projectService)
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService, projectService, staticPath)
	reportHandler := handlers.NewReportHandler(reportService, projectService)
	analyticsHandler := handlers.NewAnalyticsHandler(reportService, projectService)

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
	}
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// -------------------------
	// Protected (Authentication required)
	// -------------------------
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware)

	// ---------- Projects ----------
	projects := protected.Group("/projects")
	{
		projects.GET("", projectHandler.ListProjects)
		projects.GET("/:projectId", projectHandler.GetProjectByID)

		admin := projects.Group("")
		admin.Use(middleware.RequiredRole("admin", "superadmin"))
		{
			admin.POST("", projectHandler.CreateProject)
			admin.PUT("/:projectId", projectHandler.UpdateProject)
			admin.DELETE("/:projectId", projectHandler.DeleteProject)
		}
	}

	// ---------- Teams ----------
	teams := protected.Group("/teams")
	teams.Use(middleware.RequiredRole("superadmin"))
	{
		teams.GET("", teamHandler.ListTeams)
		teams.GET("/:teamId", teamHandler.GetTeamByID)
		teams.POST("", teamHandler.CreateTeam)
		teams.PUT("/:teamId", teamHandler.UpdateTeam)
		teams.DELETE("/:teamId", teamHandler.DeleteTeam)
	}

	// ---------- Users ----------
	users := protected.Group("/users")
	{
		users.GET("/me", userHandler.GetCurrentUser)
		users.GET("/:userID", userHandler.GetUserByID)
		users.PATCH("/:userID/password", userHandler.ChangePassword)
		users.PATCH("/:userID/project", userHandler.SetDefaultProject)
		users.PATCH("/:userID/image", userHandler.ChangeImage)

		admin := users.Group("")
		admin.Use(middleware.RequiredRole("superadmin"))
		{
			admin.GET("", userHandler.ListUsers)
			admin.PATCH("/:userID", userHandler.UpdateUser)
			admin.DELETE("/:userID", userHandler.DeleteUser)
		}
	}

	// ---------- Reports ----------
	reports := protected.Group("/reports")
	{
		reports.GET("/:projectID", reportHandler.ListReportsByProjectID)
	}
	api.POST("/reports/:projectID/endpoint", reportHandler.CreateReport)

	// ---------- Analytics ----------
	analytics := protected.Group("/analytics")
	{
		analytics.GET("/:projectID/graph-data", analyticsHandler.GetReportGraphData)
		analytics.GET("/:projectID/summary-stats", analyticsHandler.GetReportSummaryStats)
		analytics.GET("/:projectID/violation-trend", analyticsHandler.GetReportViolationTrend)
		analytics.GET("/:projectID/violated-directives", analyticsHandler.GetReportTopViolatedDirectives)
		analytics.GET("/:projectID/violated-document-urls", analyticsHandler.GetReportTopViolatedDocumentURLs)
		analytics.GET("/:projectID/software-stats", analyticsHandler.GetReportSoftwareStats)
		analytics.GET("/:projectID/violation-sources", analyticsHandler.GetReportTopViolationSources)
	}

	// ---------- Signout ----------
	protected.POST("/signout", authHandler.SignOut)

	return router
}
