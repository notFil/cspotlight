package services

import (
	"errors"

	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/internal/repositories"
	"github.com/notFil/cspotlight/pkg/auth"
)

type ReportService interface {
	CreateReport(report *models.CSPReportCreateDTO, projectID string) error
	DeleteReport(id string) error
	ListReportsByProjectID(projectID string, p *pagination.Pagination, claims auth.Claims) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(reports []*models.CSPReportCreateDTO, projectID string) error
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

func (s *reportService) CreateReport(report *models.CSPReportCreateDTO, projectID string) error {
	r := report.ToCSPReport()
	r.ProjectID = projectID
	return s.reportRepo.CreateReport(r)
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

func (s *reportService) DeleteReport(id string) error {
	return s.reportRepo.DeleteReport(id)
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
