package services

import (
	"context"
	"net/http"

	apperrors "cspotlight/internal/errors"
	"cspotlight/internal/models"
	"cspotlight/internal/pagination"
	"cspotlight/internal/repositories"

	"github.com/google/uuid"
)

type ReportService interface {
	ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID uuid.UUID) error
	GetReportSummaryStats(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error)
	GetReportGraphData(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error)
	GetReportViolationTrend(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error)
	GetReportTopViolatedDocumentURLs(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error)
	GetReportTopViolatedDirectives(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error)
	GetReportSoftwareStats(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error)
	GetReportTopViolationSources(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}

func (s *reportService) BatchCreateReports(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID uuid.UUID) error {
	var cspReports []*models.CSPReport
	for _, r := range reports {
		cspReport := r.ToCSPReport()
		cspReport.ProjectID = projectID
		cspReports = append(cspReports, cspReport)
	}

	return s.reportRepo.BatchCreateReports(ctx, cspReports)
}

func (s *reportService) ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	reports, p, err := s.reportRepo.ListReportsByProjectID(ctx, projectID, p)
	if err != nil {
		return nil, p, err
	}

	return reports, p, nil
}

func (s *reportService) GetReportSummaryStats(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error) {
	stats, err := s.reportRepo.GetReportSummaryStats(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, apperrors.New(http.StatusNotFound, "report summary stats not found")
	}

	return stats, nil
}

func (s *reportService) GetReportGraphData(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error) {
	g, err := s.reportRepo.GetReportGraphData(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if g == nil {
		return nil, apperrors.New(http.StatusNotFound, "report graph data not found")
	}
	return g, nil
}

func (s *reportService) GetReportViolationTrend(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error) {
	t, err := s.reportRepo.GetReportViolationTrend(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, apperrors.New(http.StatusNotFound, "report not found")
	}
	return t, nil
}

func (s *reportService) GetReportTopViolatedDirectives(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error) {
	t, err := s.reportRepo.GetReportTopViolatedDirectives(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, apperrors.New(http.StatusNotFound, "report not found")
	}
	return t, nil
}

func (s *reportService) GetReportTopViolatedDocumentURLs(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error) {
	t, err := s.reportRepo.GetReportTopViolatedDocumentURLs(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, apperrors.New(http.StatusNotFound, "report not found")
	}
	return t, nil
}

func (s *reportService) GetReportSoftwareStats(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error) {
	stats, err := s.reportRepo.GetReportSoftwareStats(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, apperrors.New(http.StatusNotFound, "report software stats not found")
	}
	return stats, nil
}

func (s *reportService) GetReportTopViolationSources(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error) {
	t, err := s.reportRepo.GetReportTopViolationSources(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if t == nil {
		return nil, apperrors.New(http.StatusNotFound, "report top violation sources not found")
	}
	return t, nil
}
