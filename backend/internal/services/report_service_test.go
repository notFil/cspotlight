package services

import (
	"context"
	"errors"
	"testing"

	"cspotlight/internal/auth"
	"cspotlight/internal/models"
	"cspotlight/internal/pagination"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// MockReportRepository is a manual mock for ReportRepository
type MockReportRepository struct {
	reports    map[uuid.UUID]*models.CSPReport
	reportDTOs map[uuid.UUID]*models.CSPReportFetchDTO
}

func NewMockReportRepository() *MockReportRepository {
	return &MockReportRepository{
		reports:    make(map[uuid.UUID]*models.CSPReport),
		reportDTOs: make(map[uuid.UUID]*models.CSPReportFetchDTO),
	}
}

func (m *MockReportRepository) CreateReport(ctx context.Context, report *models.CSPReport) error {
	if report.ID == uuid.Nil {
		report.ID = uuid.New()
	}
	m.reports[report.ID] = report
	return nil
}

func (m *MockReportRepository) BatchCreateReports(ctx context.Context, reports []*models.CSPReport) error {
	for _, report := range reports {
		if report.ID == uuid.Nil {
			report.ID = uuid.New()
		}
		m.reports[report.ID] = report
	}
	return nil
}

func (m *MockReportRepository) GetReportByID(ctx context.Context, id uuid.UUID) (*models.CSPReport, error) {
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

func (m *MockReportRepository) DeleteReport(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.reports[id]; exists {
		delete(m.reports, id)
		return nil
	}
	return errors.New("report not found")
}

func (m *MockReportRepository) ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination, filter *models.ReportFilter) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
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

	p.TotalRows = int64(len(reports))
	if p.PageSize <= 0 {
		p.PageSize = 50
	}
	p.TotalPages = int((p.TotalRows + int64(p.PageSize) - 1) / int64(p.PageSize))

	return reports, p, nil
}

func (m *MockReportRepository) GetReportSummaryStats(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error) {
	return &models.ReportMetricsDTO{}, nil
}

func (m *MockReportRepository) GetReportGraphData(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error) {
	return &models.ReportGraphDataDTO{}, nil
}

func (m *MockReportRepository) GetReportViolationTrend(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error) {
	return &models.ReportViolationTrendDTO{}, nil
}

func (m *MockReportRepository) GetReportTopViolatedDocumentURLs(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error) {
	return &models.ReportTopViolatedDocumentURLDTO{}, nil
}

func (m *MockReportRepository) GetReportTopViolatedDirectives(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error) {
	return &models.ReportTopViolatedDirectivesDTO{}, nil
}

func (m *MockReportRepository) GetReportSoftwareStats(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error) {
	return &models.ReportSoftwareStatsDTO{}, nil
}

func (m *MockReportRepository) GetReportTopViolationSources(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error) {
	return &models.ReportTopViolationSourcesDTO{}, nil
}

func TestCreateReport(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()

	service := NewReportService(mockRepo, mockProjectRepo)

	projectID := uuid.New()
	mockProjectRepo.projects[projectID] = &models.Project{ID: projectID, Disabled: false}
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

	projectID := uuid.New()
	teamID := uuid.New()

	mockProjectRepo.projects[projectID] = &models.Project{ID: projectID, TeamID: teamID}

	r1 := &models.CSPReport{ID: uuid.New(), ProjectID: projectID, Body: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b1"})}
	r2 := &models.CSPReport{ID: uuid.New(), ProjectID: projectID, Body: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b2"})}
	r3 := &models.CSPReport{ID: uuid.New(), ProjectID: uuid.New(), Body: datatypes.NewJSONType(models.ReportBody{BlockedURL: "b3"})}

	mockRepo.reports[r1.ID] = r1
	mockRepo.reports[r2.ID] = r2
	mockRepo.reports[r3.ID] = r3

	p := &pagination.Pagination{Page: 1, PageSize: 10}
	userData := auth.UserContext{ID: uuid.New(), Role: "user", TeamID: teamID}
	ctx := auth.ContextWithUser(context.Background(), &userData)

	dtos, _, err := service.ListReportsByProjectID(ctx, projectID, p, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(dtos) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(dtos))
	}
}

func TestGetReportSummaryStats(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	stats, err := service.GetReportSummaryStats(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats == nil {
		t.Fatal("expected stats, got nil")
	}
}

func TestGetReportGraphData(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	data, err := service.GetReportGraphData(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatal("expected data, got nil")
	}
}

func TestGetReportViolationTrend(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	data, err := service.GetReportViolationTrend(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatal("expected data, got nil")
	}
}

func TestGetReportTopViolatedDocumentURLs(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	data, err := service.GetReportTopViolatedDocumentURLs(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatal("expected data, got nil")
	}
}

func TestGetReportTopViolatedDirectives(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	data, err := service.GetReportTopViolatedDirectives(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatal("expected data, got nil")
	}
}

func TestGetReportSoftwareStats(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	data, err := service.GetReportSoftwareStats(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatal("expected data, got nil")
	}
}

func TestGetReportTopViolationSources(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)
	projectID := uuid.New()

	data, err := service.GetReportTopViolationSources(context.Background(), projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data == nil {
		t.Fatal("expected data, got nil")
	}
}

func TestBatchCreateReports_ProjectStatus(t *testing.T) {
	mockRepo := NewMockReportRepository()
	mockProjectRepo := NewMockProjectRepository()
	service := NewReportService(mockRepo, mockProjectRepo)

	projectID := uuid.New()
	reportDTO := &models.CSPReportCreateDTO{
		URL: "https://example.com",
	}

	err := service.BatchCreateReports(context.Background(), []*models.CSPReportCreateDTO{reportDTO}, projectID)
	if err == nil {
		t.Fatal("expected error for non-existent project, got nil")
	}

	mockProjectRepo.projects[projectID] = &models.Project{ID: projectID, Disabled: true}
	err = service.BatchCreateReports(context.Background(), []*models.CSPReportCreateDTO{reportDTO}, projectID)
	if err != nil {
		t.Fatalf("expected no error for disabled project (skipping), got %v", err)
	}
	if len(mockRepo.reports) != 0 {
		t.Errorf("expected 0 reports for disabled project, got %d", len(mockRepo.reports))
	}

	mockProjectRepo.projects[projectID].Disabled = false
	err = service.BatchCreateReports(context.Background(), []*models.CSPReportCreateDTO{reportDTO}, projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(mockRepo.reports) != 1 {
		t.Errorf("expected 1 report for enabled project, got %d", len(mockRepo.reports))
	}
}
