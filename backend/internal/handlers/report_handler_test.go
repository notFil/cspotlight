package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/pkg/auth"
)

type MockReportService struct {
	BatchCreateReportsFunc func(reports []*models.CSPReportCreateDTO, projectID string) error
	mu                     sync.Mutex
	batchCalls             int
	receivedReports        []*models.CSPReportCreateDTO
}

func (m *MockReportService) CreateReport(report *models.CSPReportCreateDTO, projectID string) error {
	return nil
}

func (m *MockReportService) GetReportByID(id string) (*models.CSPReport, error) {
	return nil, nil
}

func (m *MockReportService) UpdateReport(report *models.CSPReport) error {
	return nil
}

func (m *MockReportService) DeleteReport(id string) error {
	return nil
}

func (m *MockReportService) ListReportsByProjectID(projectID string, p *pagination.Pagination, claims auth.Claims) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	return nil, nil, nil
}

func (m *MockReportService) BatchCreateReports(reports []*models.CSPReportCreateDTO, projectID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.batchCalls++

	// Verify SourceIP is set
	m.receivedReports = append(m.receivedReports, reports...)

	for _, r := range reports {
		if r.SourceIP != "1.2.3.4" {
			// We can't easily fail the test from here without passing *testing.T
			// But we can panic or log. For now let's just print.
			// Better: store it and check in test function.
		}
	}

	if m.BatchCreateReportsFunc != nil {
		return m.BatchCreateReportsFunc(reports, projectID)
	}
	return nil
}

func TestReportHandler_CreateReport_Batching(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockReportService{}
	handler := NewReportHandler(mockService)

	// Send 105 reports. Batch size is 100.
	// We expect 1 batch call immediately after 100, and another one after timeout for the remaining 5.
	// Or just 1 batch call if we check quickly.

	// Actually, the worker runs in background.
	// Let's send 150 reports.
	totalReports := 150

	for i := 0; i < totalReports; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "projectID", Value: "proj-1"}}

		body := `{"reportBody": {"blockedURL": "http://evil.com"}, "type": "csp-report"}`
		c.Request = httptest.NewRequest("POST", "/reports/proj-1", strings.NewReader(body))
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
