package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/middleware"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
)

type MockReportService struct {
	BatchCreateReportsFunc               func(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID uuid.UUID) error
	ListReportsByProjectIDFunc           func(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	GetReportSummaryStatsFunc            func(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error)
	GetReportGraphDataFunc               func(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error)
	GetReportViolationTrendFunc          func(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error)
	GetReportTopViolatedDocumentURLsFunc func(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error)
	GetReportTopViolatedDirectivesFunc   func(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error)
	GetReportSoftwareStatsFunc           func(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error)
	GetReportTopViolationSourcesFunc     func(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error)
	mu                                   sync.Mutex // Kept for potential future use or if other tests rely on it
	batchCalls                           int        // Kept for potential future use or if other tests rely on it
	receivedReports                      []*models.CSPReportCreateDTO
}

func (m *MockReportService) ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	if m.ListReportsByProjectIDFunc != nil {
		return m.ListReportsByProjectIDFunc(ctx, projectID, p)
	}
	return nil, nil, nil
}

func (m *MockReportService) BatchCreateReports(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.batchCalls++

	// Capture reports for test verification
	m.receivedReports = append(m.receivedReports, reports...)

	if m.BatchCreateReportsFunc != nil {
		return m.BatchCreateReportsFunc(ctx, reports, projectID)
	}
	return nil
}

func (m *MockReportService) GetReportGraphData(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error) {
	return nil, nil
}

func (m *MockReportService) GetReportSummaryStats(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error) {
	return nil, nil
}

func (m *MockReportService) GetReportViolationTrend(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error) {
	return nil, nil
}

func (m *MockReportService) GetReportTopViolatedDirectives(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error) {
	return nil, nil
}

func (m *MockReportService) GetReportSoftwareStats(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error) {
	return nil, nil
}

func (m *MockReportService) GetReportTopViolatedDocumentURLs(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error) {
	return nil, nil
}

func (m *MockReportService) GetReportTopViolationSources(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error) {
	return nil, nil
}

func TestReportHandler_CreateReport_Batching(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockReportService{}
	mockProjectService := &MockProjectService{}
	handler := NewReportHandler(mockService, mockProjectService)

	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.POST("/reports/:projectID", handler.CreateReport)

	totalReports := 150

	for i := 0; i < totalReports; i++ {
		w := httptest.NewRecorder()
		body := `{"reportBody": {"blockedURL": "http://evil.com"}, "type": "csp-report"}`
		req := httptest.NewRequest("POST", "/reports/00000000-0000-0000-0000-000000000001", strings.NewReader(body))
		req.RemoteAddr = "1.2.3.4:1234"
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("iteration %d: expected status 204, got %d", i, w.Code)
		}
	}

	// Wait for batch processing
	time.Sleep(100 * time.Millisecond)

	mockService.mu.Lock()
	calls := mockService.batchCalls
	received := mockService.receivedReports
	mockService.mu.Unlock()

	if calls < 1 {
		t.Errorf("expected at least 1 batch call, got %d", calls)
	}

	if len(received) == 0 {
		t.Error("expected received reports, got 0")
	}

	for _, r := range received {
		if r.SourceIP != "1.2.3.4" {
			t.Errorf("expected SourceIP 1.2.3.4, got %s", r.SourceIP)
		}
	}
}

func TestReportHandler_ListReportsByProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService := &MockReportService{
			ListReportsByProjectIDFunc: func(ctx context.Context, pid uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
				return []*models.CSPReportFetchDTO{{BlockedURL: "http://example.com"}}, &pagination.Pagination{Page: 1, PageSize: 10, TotalRows: 1}, nil
			},
		}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return &models.ProjectFetchDTO{ID: id}, nil
			},
		}
		handler := NewReportHandler(mockService, mockProjectService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/reports/:projectID", handler.ListReportsByProjectID)

		req := httptest.NewRequest("GET", "/reports/"+projectID.String(), nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("ProjectNotFound", func(t *testing.T) {
		mockService := &MockReportService{}
		mockProjectService := &MockProjectService{
			GetProjectByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.ProjectFetchDTO, error) {
				return nil, apperrors.New(http.StatusNotFound, "project not found")
			},
		}
		handler := NewReportHandler(mockService, mockProjectService)

		w := httptest.NewRecorder()
		router := gin.New()
		router.Use(middleware.ErrorHandler())
		router.GET("/reports/:projectID", handler.ListReportsByProjectID)

		req := httptest.NewRequest("GET", "/reports/"+projectID.String(), nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}
