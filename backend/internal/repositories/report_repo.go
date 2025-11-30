package repositories

import (
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"gorm.io/gorm"
)

type ReportRepository interface {
	CreateReport(report *models.CSPReport) error
	GetReportByID(id string) (*models.CSPReport, error)
	UpdateReport(report *models.CSPReport) error
	DeleteReport(id string) error
	ListReportsByProjectID(projectID string, p *pagination.Pagination) ([]*models.CSPReport, *pagination.Pagination, error)
	BatchCreateReports(reports []*models.CSPReport) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) CreateReport(report *models.CSPReport) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) BatchCreateReports(reports []*models.CSPReport) error {
	return r.db.CreateInBatches(reports, 100).Error
}

func (r *reportRepository) GetReportByID(id string) (report *models.CSPReport, err error) {
	err = r.db.First(&report, "id = ?", id).Error
	return report, err
}

func (r *reportRepository) UpdateReport(report *models.CSPReport) error {
	return r.db.Save(report).Error
}

func (r *reportRepository) DeleteReport(id string) error {
	return r.db.Delete(&models.CSPReport{}, "id = ?", id).Error
}

func (r *reportRepository) ListReportsByProjectID(projectID string, p *pagination.Pagination) ([]*models.CSPReport, *pagination.Pagination, error) {
	var reports []*models.CSPReport

	r.db.Where("project_id = ?", projectID).Count(&p.TotalRows)

	p.TotalPages = int((p.TotalRows + int64(p.PageSize) - 1) / int64(p.PageSize))

	result := r.db.Scopes(pagination.Paginate(p)).Order("id DESC").Find(&reports)
	return reports, p, result.Error
}
