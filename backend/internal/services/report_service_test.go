package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"gorm.io/datatypes"
)

// MockReportRepository is a manual mock for ReportRepository
type MockReportRepository struct {
	reports    map[string]*models.CSPReport
	reportDTOs map[string]*models.CSPReportFetchDTO
}

func NewMockReportRepository() *MockReportRepository {
	return &MockReportRepository{
		reports:    make(map[string]*models.CSPReport),
		reportDTOs: make(map[string]*models.CSPReportFetchDTO),
	}
}

func (m *MockReportRepository) CreateReport(ctx context.Context, report *models.CSPReport) error {
	if report.ID == "" {
		report.ID = uuid.New().String()
	}
	m.reports[report.ID] = report
	return nil
}

func (m *MockReportRepository) BatchCreateReports(ctx context.Context, reports []*models.CSPReport) error {
	for _, report := range reports {
		if report.ID == "" {
			report.ID = uuid.New().String()
		}
		m.reports[report.ID] = report
	}
	return nil
}

func (m *MockReportRepository) GetReportByID(ctx context.Context, id string) (*models.CSPReport, error) {
	if report, exists := m.reports[id]; exists {
		return report, nil
	}
	return nil, errors.New("report not found")
}

func (m *MockReportRepository) UpdateReport(ctx context.Context, report *models.CSPReport) error {
	if _, exists := m.reports[report.ID]; exists {
		m.reports[report.ID] = report
		return nil
	}
	return errors.New("report not found")
}

func (m *MockReportRepository) DeleteReport(ctx context.Context, id string) error {
	if _, exists := m.reports[id]; exists {
		delete(m.reports, id)
		return nil
	}
	return errors.New("report not found")
}

func (m *MockReportRepository) ListReportsByProjectID(ctx context.Context, projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	var reports []*models.CSPReportFetchDTO
	for _, report := range m.reports {
		if report.ProjectID == projectID {
			dto := &models.CSPReportFetchDTO{
				Body:        report.Body,
				Directive:   report.Directive,
				URL:         report.URL,
				BlockedURL:  report.BlockedURL,
				DocumentURL: report.DocumentURL,
				Disposition: report.Disposition,
				SourceIP:    report.SourceIP,
				UserAgent:   report.UserAgent,
				LastSeen:    report.CreatedAt,
				Count:       1,
			}
			reports = append(reports, dto)
		}
	}

	// Simple pagination logic for mock
	p.TotalRows = int64(len(reports))
	if p.PageSize <= 0 {
		p.PageSize = 50
	}
	p.TotalPages = int((p.TotalRows + int64(p.PageSize) - 1) / int64(p.PageSize))

	return reports, p, nil
}

func (m *MockReportRepository) GetReportSummaryStats(ctx context.Context, projectID string) (*models.ReportMetricsDTO, error) {
	return &models.ReportMetricsDTO{}, nil
}

func (m *MockReportRepository) GetReportGraphData(ctx context.Context, projectID string) (*models.ReportGraphDataDTO, error) {
	return &models.ReportGraphDataDTO{}, nil
}

func TestCreateReport(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	projectID := "proj-1"
	reportDTO := &models.CSPReportCreateDTO{
		Age:       100,
		Type:      "csp-report",
		URL:       "https://example.com",
		UserAgent: "Mozilla/5.0",
		Body: models.ReportBody{
			BlockedURL: "https://evil.com",
		},
	}

	err := service.BatchCreateReports(context.Background(), []*models.CSPReportCreateDTO{reportDTO}, projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify report was created
	if len(mockRepo.reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(mockRepo.reports))
	}

	var createdReport *models.CSPReport
	for _, r := range mockRepo.reports {
		createdReport = r
		break
	}

	if createdReport.ProjectID != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, createdReport.ProjectID)
	}
	if createdReport.URL != reportDTO.URL {
		t.Errorf("expected URL %s, got %s", reportDTO.URL, createdReport.URL)
	}
}

func TestListReportsByProjectID(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	projectID := "proj-1"
	teamID := "team-1"

	// Setup project
	mockProjectRepo.projects[projectID] = &models.Project{ID: projectID, TeamID: teamID}

	// Add reports
	r1 := &models.CSPReport{ID: "r1", ProjectID: projectID, Body: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b1"})}
	r2 := &models.CSPReport{ID: "r2", ProjectID: projectID, Body: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b2"})}
	r3 := &models.CSPReport{ID: "r3", ProjectID: "other-proj", Body: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b3"})}

	mockRepo.reports["r1"] = r1
	mockRepo.reports["r2"] = r2
	mockRepo.reports["r3"] = r3

	p := &pagination.Pagination{Page: 1, PageSize: 10}
	claims := auth.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "user"}, Role: "user", TeamID: teamID}
	ctx := auth.ContextWithClaims(context.Background(), &claims)

	dtos, _, err := service.ListReportsByProjectID(ctx, projectID, p)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(dtos) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(dtos))
	}
}
