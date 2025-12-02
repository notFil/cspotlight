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
	ListReportsByProjectID(projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
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

func (r *reportRepository) ListReportsByProjectID(projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	var reports []*models.CSPReportFetchDTO

	// Count total unique groups for pagination
	var totalRows int64
	r.db.Raw(`
		SELECT COUNT(*) 
		FROM (
			SELECT 1 
			FROM csp_reports 
			WHERE project_id = ? 
			GROUP BY url, directive, blocked_url, disposition, document_url, body, source_ip, user_agent
		) AS sub
	`, projectID).Scan(&totalRows)
	p.TotalRows = totalRows

	p.TotalPages = int((p.TotalRows + int64(p.GetLimit()) - 1) / int64(p.GetLimit()))

	// Fetch paginated results
	result := r.db.Raw(`
    SELECT 
        url,
        directive,
				blocked_url,
				disposition,
				document_url,
				body,
				source_ip,
				user_agent,
        COUNT(*) AS count,
        MAX(created_at) AS last_seen
    FROM csp_reports
    WHERE project_id = ?
    GROUP BY url, directive, blocked_url, disposition, document_url, body, source_ip, user_agent
    LIMIT ? OFFSET ?
`, projectID, p.GetLimit(), p.GetOffset()).Scan(&reports)

	if result.Error != nil {
		return nil, p, result.Error
	}

	return reports, p, nil
}
