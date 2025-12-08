package services

import (
	"context"
	"net/http"

	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/internal/repositories"
)

type ReportService interface {
	ListReportsByProjectID(ctx context.Context, projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID string) error
	GetReportSummaryStats(ctx context.Context, projectID string) (*models.ReportMetricsDTO, error)
	GetReportGraphData(ctx context.Context, projectID string) (*models.ReportGraphDataDTO, error)
	GetReportViolationTrend(ctx context.Context, projectID string) (*models.ReportViolationTrendDTO, error)
}

type reportService struct {
	reportRepo  repositories.ReportRepository
	projectRepo repositories.ProjectRepository
}

func NewReportService(reportRepo repositories.ReportRepository, projectRepo repositories.ProjectRepository) ReportService {
	return &reportService{
		reportRepo:  reportRepo,
		projectRepo: projectRepo,
	}
}

func (s *reportService) BatchCreateReports(ctx context.Context, reports []*models.CSPReportCreateDTO, projectID string) error {
	var cspReports []*models.CSPReport
	for _, r := range reports {
		cspReport := r.ToCSPReport()
		cspReport.ProjectID = projectID
		cspReports = append(cspReports, cspReport)
	}

	return s.reportRepo.BatchCreateReports(ctx, cspReports)
}

func (s *reportService) ListReportsByProjectID(ctx context.Context, projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	claims := auth.GetUserClaims(ctx)
	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, p, err
	}
	if claims.IsSuperadmin() && project.TeamID != claims.TeamID {
		return nil, p, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	reports, p, err := s.reportRepo.ListReportsByProjectID(ctx, projectID, p)
	if err != nil {
		return nil, p, err
	}

	return reports, p, nil
}

func (s *reportService) GetReportSummaryStats(ctx context.Context, projectID string) (*models.ReportMetricsDTO, error) {
	stats := &models.ReportMetricsDTO{}

	claims := auth.GetUserClaims(ctx)

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && project.TeamID != claims.TeamID {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if stats, err = s.reportRepo.GetReportSummaryStats(ctx, projectID); err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, apperrors.New(http.StatusNotFound, "report summary stats not found")
	}

	return stats, nil
}

func (s *reportService) GetReportGraphData(ctx context.Context, projectID string) (*models.ReportGraphDataDTO, error) {
	g := &models.ReportGraphDataDTO{}

	claims := auth.GetUserClaims(ctx)

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && project.TeamID != claims.TeamID {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if g, err = s.reportRepo.GetReportGraphData(ctx, projectID); err != nil {
		return nil, err
	}

	if g == nil {
		return nil, apperrors.New(http.StatusNotFound, "report graph data not found")
	}
	return g, nil
}

func (s *reportService) GetReportViolationTrend(ctx context.Context, projectID string) (*models.ReportViolationTrendDTO, error) {
	t := &models.ReportViolationTrendDTO{}

	claims := auth.GetUserClaims(ctx)

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !claims.IsSuperadmin() && project.TeamID != claims.TeamID {
		return nil, apperrors.New(http.StatusUnauthorized, "unauthorized access")
	}

	if t, err = s.reportRepo.GetReportViolationTrend(ctx, projectID); err != nil {
		return nil, err
	}

	if t == nil {
		return nil, apperrors.New(http.StatusNotFound, "report not found")
	}
	return t, nil
}
