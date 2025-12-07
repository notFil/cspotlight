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
	GetReportGraphData(ctx context.Context, projectID string) (*models.ReportGraphDataDTO, error)
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
	return g, nil
}
