package services

import (
	"errors"

	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/internal/repositories"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/errs"
)

type ReportService interface {
	ListReportsByProjectID(projectID string, p *pagination.Pagination, claims auth.Claims) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(reports []*models.CSPReportCreateDTO, projectID string) error
	GetReportGraphData(projectID string, claims auth.Claims) (*models.ReportGraphDataDTO, error)
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

func (s *reportService) BatchCreateReports(reports []*models.CSPReportCreateDTO, projectID string) error {
	var cspReports []*models.CSPReport
	for _, r := range reports {
		cspReport := r.ToCSPReport()
		cspReport.ProjectID = projectID
		cspReports = append(cspReports, cspReport)
	}
	return s.reportRepo.BatchCreateReports(cspReports)
}
func (s *reportService) ListReportsByProjectID(projectID string, p *pagination.Pagination, claims auth.Claims) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	project, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, p, err
	}
	if claims.Role != "superadmin" && project.TeamID != claims.TeamID {
		return nil, p, errors.New("unauthorized access")
	}

	reports, p, err := s.reportRepo.ListReportsByProjectID(projectID, p)
	if err != nil {
		return nil, p, err
	}

	return reports, p, nil
}

func (s *reportService) GetReportGraphData(projectID string, claims auth.Claims) (*models.ReportGraphDataDTO, error) {
	project, err := s.projectRepo.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}
	if claims.Role != "superadmin" && project.TeamID != claims.TeamID {
		return nil, errs.ErrUnauthorized
	}
	return s.reportRepo.GetReportGraphData(projectID)
}
