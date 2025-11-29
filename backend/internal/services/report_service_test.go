package services

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/pkg/auth"
	"gorm.io/datatypes"
)

// MockReportRepository is a manual mock for ReportRepository
type MockReportRepository struct {
	reports map[string]*models.CSPReport
}

func NewMockReportRepository() *MockReportRepository {
	return &MockReportRepository{
		reports: make(map[string]*models.CSPReport),
	}
}

func (m *MockReportRepository) CreateReport(report *models.CSPReport) error {
	if report.ID == "" {
		report.ID = uuid.New().String()
	}
	m.reports[report.ID] = report
	return nil
}

func (m *MockReportRepository) GetReportByID(id string) (*models.CSPReport, error) {
	if report, exists := m.reports[id]; exists {
		return report, nil
	}
	return nil, errors.New("report not found")
}

func (m *MockReportRepository) UpdateReport(report *models.CSPReport) error {
	if _, exists := m.reports[report.ID]; exists {
		m.reports[report.ID] = report
		return nil
	}
	return errors.New("report not found")
}

func (m *MockReportRepository) DeleteReport(id string) error {
	if _, exists := m.reports[id]; exists {
		delete(m.reports, id)
		return nil
	}
	return errors.New("report not found")
}

func (m *MockReportRepository) ListReportsByProjectID(projectID string, p *pagination.Pagination) ([]*models.CSPReport, *pagination.Pagination, error) {
	var reports []*models.CSPReport
	for _, report := range m.reports {
		if report.ProjectID == projectID {
			reports = append(reports, report)
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
		ReportBody: models.ReportBody{
			BlockedURL: "https://evil.com",
		},
	}

	err := service.CreateReport(reportDTO, projectID)
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

func TestGetReportByID(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	reportID := "rep-1"
	report := &models.CSPReport{
		ID:        reportID,
		ProjectID: "proj-1",
		URL:       "https://example.com",
	}
	mockRepo.reports[reportID] = report

	// Test success
	fetchedReport, err := service.GetReportByID(reportID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fetchedReport.ID != reportID {
		t.Errorf("expected report ID %s, got %s", reportID, fetchedReport.ID)
	}

	// Test not found
	_, err = service.GetReportByID("non-existent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateReport(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	reportID := "rep-update"
	report := &models.CSPReport{
		ID:        reportID,
		ProjectID: "proj-1",
		URL:       "https://example.com",
	}
	mockRepo.reports[reportID] = report

	report.URL = "https://updated.com"
	err := service.UpdateReport(report)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mockRepo.reports[reportID].URL != "https://updated.com" {
		t.Errorf("expected updated URL, got %s", mockRepo.reports[reportID].URL)
	}
}

func TestDeleteReport(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	reportID := "rep-delete"
	report := &models.CSPReport{
		ID:        reportID,
		ProjectID: "proj-1",
	}
	mockRepo.reports[reportID] = report

	err := service.DeleteReport(reportID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, exists := mockRepo.reports[reportID]; exists {
		t.Fatal("expected report to be deleted")
	}
}

func TestListReportsByProjectID(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	projectID := "proj-1"
	otherProjectID := "proj-2"
	teamID := "team-1"

	// Setup project
	mockProjectRepo.projects[projectID] = &models.Project{ID: projectID, TeamID: teamID}

	// Add reports
	r1 := &models.CSPReport{ID: "r1", ProjectID: projectID, ReportBody: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b1"})}
	r2 := &models.CSPReport{ID: "r2", ProjectID: projectID, ReportBody: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b2"})}
	r3 := &models.CSPReport{ID: "r3", ProjectID: otherProjectID, ReportBody: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b3"})}

	mockRepo.reports["r1"] = r1
	mockRepo.reports["r2"] = r2
	mockRepo.reports["r3"] = r3

	p := &pagination.Pagination{Page: 1, PageSize: 10}
	claims := auth.Claims{TeamID: teamID}

	dtos, _, err := service.ListReportsByProjectID(projectID, p, claims)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(dtos) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(dtos))
	}
}
