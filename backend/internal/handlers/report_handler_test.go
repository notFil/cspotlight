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
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
)

type MockReportService struct {
	BatchCreateReportsFunc func(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID uuid.UUID) error
	mu                     sync.Mutex
	batchCalls             int
	receivedReports        []*models.CSPReportCreateDTO
}

func (m *MockReportService) ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	return nil, nil, nil
}

func (m *MockReportService) BatchCreateReports(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.batchCalls++

	// Verify SourceIP is set
	m.receivedReports = append(m.receivedReports, reports...)

	for _, r := range reports {
		if r.SourceIP != "1.2.3.4" {
			// logging or panic?
		}
	}

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
	handler := NewReportHandler(mockService)

	totalReports := 150

	for i := 0; i < totalReports; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: "00000000-0000-0000-0000-000000000001"}}

		body := `{"reportBody": {"blockedURL": "http://evil.com"}, "type": "csp-report"}`
		c.Request = httptest.NewRequest("POST", "/reports/00000000-0000-0000-0000-000000000001", strings.NewReader(body))
		c.Request.RemoteAddr = "1.2.3.4:1234"

		handler.CreateReport(c)

		if c.Writer.Status() != http.StatusNoContent {
			t.Errorf("iteration %d: expected status 204, got %d", i, c.Writer.Status())
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
