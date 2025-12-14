package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/middleware"
	"github.com/notFil/cspotlight/internal/models"
)

func TestAnalyticsHandler_GetReportSummaryStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportSummaryStatsFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportMetricsDTO, error) {
				return &models.ReportMetricsDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/analytics/:projectID/summary-stats", handler.GetReportSummaryStats)

		req := httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/summary-stats", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("ProjectNotFound", func(t *testing.T) {
		mockReportService := &MockReportService{}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return nil, apperrors.New(http.StatusNotFound, "project not found")
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/analytics/:projectID/summary-stats", handler.GetReportSummaryStats)

		req := httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/summary-stats", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestAnalyticsHandler_GetReportGraphData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportGraphDataFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportGraphDataDTO, error) {
				return &models.ReportGraphDataDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: projectID.String()}}
		c.Request = httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/graph-data", nil)

		handler.GetReportGraphData(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAnalyticsHandler_GetReportTopViolatedDirectives(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportTopViolatedDirectivesFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error) {
				return &models.ReportTopViolatedDirectivesDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: projectID.String()}}
		c.Request = httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/top-violated-directives", nil)

		handler.GetReportTopViolatedDirectives(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAnalyticsHandler_GetReportViolationTrend(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportViolationTrendFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportViolationTrendDTO, error) {
				return &models.ReportViolationTrendDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: projectID.String()}}
		c.Request = httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/violation-trend", nil)

		handler.GetReportViolationTrend(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAnalyticsHandler_GetReportTopViolatedDocumentURLs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportTopViolatedDocumentURLsFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error) {
				return &models.ReportTopViolatedDocumentURLDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: projectID.String()}}
		c.Request = httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/top-violated-document-urls", nil)

		handler.GetReportTopViolatedDocumentURLs(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAnalyticsHandler_GetReportSoftwareStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportSoftwareStatsFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportSoftwareStatsDTO, error) {
				return &models.ReportSoftwareStatsDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: projectID.String()}}
		c.Request = httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/software-stats", nil)

		handler.GetReportSoftwareStats(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}

func TestAnalyticsHandler_GetReportTopViolationSources(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockReportService := &MockReportService{
			GetReportTopViolationSourcesFunc: func(ctx context.Context, pid uuid.UUID) (*models.ReportTopViolationSourcesDTO, error) {
				return &models.ReportTopViolationSourcesDTO{}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewAnalyticsHandler(mockReportService, mockProjectService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: projectID.String()}}
		c.Request = httptest.NewRequest("GET", "/analytics/"+projectID.String()+"/top-violation-sources", nil)

		handler.GetReportTopViolationSources(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})
}
