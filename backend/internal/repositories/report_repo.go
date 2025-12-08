package repositories

import (
	"context"

	"github.com/notFil/cspotlight/internal/constants"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"gorm.io/gorm"
)

type ReportRepository interface {
	ListReportsByProjectID(ctx context.Context, projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(ctx context.Context, reports []*models.CSPReport) error
	GetReportSummaryStats(ctx context.Context, projectID string) (*models.ReportMetricsDTO, error)
	GetReportGraphData(ctx context.Context, projectID string) (*models.ReportGraphDataDTO, error)
	GetReportViolationTrend(ctx context.Context, projectID string) (*models.ReportViolationTrendDTO, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) BatchCreateReports(ctx context.Context, reports []*models.CSPReport) error {
	return r.db.WithContext(ctx).CreateInBatches(reports, 100).Error
}

func (r *reportRepository) ListReportsByProjectID(ctx context.Context, projectID string, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	var reports []*models.CSPReportFetchDTO

	// Count total unique groups for pagination
	var totalRows int64
	r.db.WithContext(ctx).Raw(`
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
	result := r.db.WithContext(ctx).Raw(`
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

func (r *reportRepository) GetReportSummaryStats(ctx context.Context, projectID string) (*models.ReportMetricsDTO, error) {
	var dto models.ReportMetricsDTO

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(*) AS total_violations,
			COUNT(CASE WHEN created_at < NOW() - INTERVAL '24 hours' THEN 1 END) AS total_violations_24h,
			COUNT(CASE WHEN directive in ? THEN 1 END) AS total_critical_violations,
			COUNT(CASE WHEN directive in ? AND created_at < NOW() - INTERVAL '24 hours' THEN 1 END) AS total_critical_violations_24h,
			COUNT(DISTINCT document_url) AS affected_domains,
			COUNT(DISTINCT CASE WHEN created_at < NOW() - INTERVAL '24 hours' THEN document_url END) AS affected_domains_24h,
			COUNT(CASE WHEN disposition = 'enforce' THEN 1 END) AS enforce_count,
			COUNT(CASE WHEN disposition = 'enforce' AND created_at < NOW() - INTERVAL '24 hours' THEN 1 END) AS enforce_count_24h
		FROM csp_reports
		WHERE project_id = ?
	`, constants.CriticalDirectives, constants.CriticalDirectives, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var totalViolations int
		var totalViolations24h int
		var totalCriticalViolations int
		var totalCriticalViolations24h int
		var affectedDomains int
		var affectedDomains24h int
		var enforceCount int
		var enforceCount24h int
		if err := rows.Scan(&totalViolations, &totalViolations24h, &totalCriticalViolations, &totalCriticalViolations24h, &affectedDomains, &affectedDomains24h, &enforceCount, &enforceCount24h); err != nil {
			return nil, err
		}

		var enforcePercentage int
		if totalViolations > 0 {
			enforcePercentage = (enforceCount * 100) / totalViolations
		}

		var enforcePercentage24h int
		if totalViolations24h > 0 {
			enforcePercentage24h = (enforceCount24h * 100) / totalViolations24h
		}

		dto = models.ReportMetricsDTO{
			TotalViolations:         models.MetricSummary{Value: totalViolations, Change: getPercentageChange(totalViolations24h, totalViolations)},
			TotalCriticalViolations: models.MetricSummary{Value: totalCriticalViolations, Change: getPercentageChange(totalCriticalViolations24h, totalCriticalViolations)},
			AffectedDomains:         models.MetricSummary{Value: affectedDomains, Change: getPercentageChange(affectedDomains24h, affectedDomains)},
			PolicyEnforcement:       models.MetricSummary{Value: enforcePercentage, Change: getPercentageChange(enforcePercentage24h, enforcePercentage)},
		}
	}

	return &dto, nil
}

func (r *reportRepository) GetReportGraphData(ctx context.Context, projectID string) (*models.ReportGraphDataDTO, error) {
	var dto models.ReportGraphDataDTO

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
			TO_CHAR(created_at, 'FMMonth DD') AS day,
			directive,
			COUNT(id) AS count
		FROM csp_reports
		WHERE project_id = ?
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY day, directive, DATE(created_at)
		ORDER BY DATE(created_at) ASC, directive ASC
	`, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var day string
		var directive string
		var count int64
		if err := rows.Scan(&day, &directive, &count); err != nil {
			return nil, err
		}
		if len(dto) > 0 && dto[len(dto)-1].Day == day {
			dto[len(dto)-1].Violations = append(dto[len(dto)-1].Violations, models.Violation{Directive: directive, Count: int(count)})
		} else {
			dto = append(dto, models.DataPoint{Day: day, Violations: []models.Violation{{Directive: directive, Count: int(count)}}})
		}

	}

	return &dto, nil
}

func (r *reportRepository) GetReportViolationTrend(ctx context.Context, projectID string) (*models.ReportViolationTrendDTO, error) {
	var dto models.ReportViolationTrendDTO

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT 
			TO_CHAR(created_at, 'FMMonth DD') AS day,
			COUNT(*) AS total,
			COUNT(CASE WHEN directive in ? THEN 1 END) AS critical,
			COUNT(CASE WHEN directive in ? THEN 1 END) AS high,
			COUNT(CASE WHEN directive in ? THEN 1 END) AS medium
		FROM csp_reports
		WHERE project_id = ?
		AND created_at >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY day, DATE(created_at)
		ORDER BY DATE(created_at) ASC
	`, constants.CriticalDirectives, constants.HighDirectives, constants.MediumDirectives, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var day string
		var total int
		var critical int
		var high int
		var medium int
		if err := rows.Scan(&day, &total, &critical, &high, &medium); err != nil {
			return nil, err
		}

		dto = append(dto, models.ViolationTrend{Day: day, Total: total, Critical: critical, High: high, Medium: medium})

	}

	return &dto, nil
}

func getPercentageChange(oldValue, newValue int) float64 {
	if oldValue == 0 {
		if newValue == 0 {
			return 0.0
		}
		return 100.0
	}
	change := float64(newValue-oldValue) / float64(oldValue) * 100.0
	return change
}
