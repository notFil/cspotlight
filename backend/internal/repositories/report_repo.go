package repositories

import (
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"gorm.io/gorm"
)

type ReportRepository interface {
	ListReportsByProjectID(projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(reports []*models.CSPReport) error
	GetReportGraphData(projectID string) (*models.ReportGraphDataDTO, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) BatchCreateReports(reports []*models.CSPReport) error {
	return r.db.CreateInBatches(reports, 100).Error
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

func (r *reportRepository) GetReportGraphData(projectID string) (*models.ReportGraphDataDTO, error) {
	var dataMap models.ReportGraphDataDTO

	rows, err := r.db.Raw(`
		SELECT 
			(CURRENT_DATE - DATE(created_at)) AS days_ago,
			directive,
			COUNT(id) AS count
		FROM csp_reports
		WHERE project_id = ?
		AND (CURRENT_DATE - DATE(created_at)) < 30
    AND (CURRENT_DATE - DATE(created_at)) >= 0
		GROUP BY days_ago, directive
		ORDER BY days_ago ASC, directive ASC
	`, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var daysAgo int
		var directive string
		var count int64
		if err := rows.Scan(&daysAgo, &directive, &count); err != nil {
			return nil, err
		}
		dataMap = append(dataMap, models.DataPoint{DaysAgo: daysAgo, Violations: []models.Violation{{Directive: directive, Count: int(count)}}})
	}

	return &dataMap, nil
}
