package repositories

import (
	"context"
	"math"

	"cspotlight/internal/constants"
	"cspotlight/internal/models"
	"cspotlight/internal/pagination"
	"cspotlight/internal/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository interface {
	ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error)
	BatchCreateReports(ctx context.Context, reports []*models.CSPReport) error
	GetReportSummaryStats(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error)
	GetReportGraphData(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error)
	GetReportViolationTrend(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error)
	GetReportTopViolatedDirectives(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error)
	GetReportTopViolatedDocumentURLs(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error)
	GetReportSoftwareStats(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error)
	GetReportTopViolationSources(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error)
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

func (r *reportRepository) ListReportsByProjectID(ctx context.Context, projectID uuid.UUID, p *pagination.Pagination) ([]*models.CSPReportFetchDTO, *pagination.Pagination, error) {
	var reports []*models.CSPReportFetchDTO

	// Count total unique groups for pagination
	// Count total unique groups for pagination
	var totalRows int64
	countQuery := `
		SELECT COUNT(*)
		FROM (
			SELECT 1
			FROM csp_reports
			WHERE project_id = ?
			GROUP BY url, directive, blocked_url, disposition, document_url, body, source_ip, user_agent
		) AS sub
	`
	r.db.WithContext(ctx).Raw(countQuery, projectID).Scan(&totalRows)
	p.TotalRows = totalRows

	p.TotalPages = int((p.TotalRows + int64(p.GetLimit()) - 1) / int64(p.GetLimit()))

	// Fetch paginated results
	// Fetch paginated results
	selectQuery := `
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
`
	result := r.db.WithContext(ctx).Raw(selectQuery, projectID, p.GetLimit(), p.GetOffset()).Scan(&reports)

	if result.Error != nil {
		return nil, p, result.Error
	}

	return reports, p, nil
}

func (r *reportRepository) GetReportSummaryStats(ctx context.Context, projectID uuid.UUID) (*models.ReportMetricsDTO, error) {
	var dto models.ReportMetricsDTO

	query := `
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
	`
	rows, err := r.db.WithContext(ctx).Raw(query, constants.CriticalDirectives, constants.CriticalDirectives, projectID).Rows()

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
			PolicyEnforcement:       models.MetricSummary{Value: enforcePercentage, Change: float64(enforcePercentage24h - enforcePercentage)},
		}
	}

	return &dto, nil
}

func (r *reportRepository) GetReportGraphData(ctx context.Context, projectID uuid.UUID) (*models.ReportGraphDataDTO, error) {
	var dto models.ReportGraphDataDTO

	query := `
		SELECT
			TO_CHAR(created_at, 'YYYY-MM-DD') AS date,
			directive,
			COUNT(id) AS count
		FROM csp_reports
		WHERE project_id = ?
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY date, directive, DATE(created_at)
		ORDER BY DATE(created_at) ASC, directive ASC
	`
	rows, err := r.db.WithContext(ctx).Raw(query, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var directive string
		var count int64
		if err := rows.Scan(&date, &directive, &count); err != nil {
			return nil, err
		}
		if len(dto) > 0 && dto[len(dto)-1].Date == date {
			dto[len(dto)-1].Violations = append(dto[len(dto)-1].Violations, models.Violation{Directive: directive, Count: int(count)})
		} else {
			dto = append(dto, models.DataPoint{Date: date, Violations: []models.Violation{{Directive: directive, Count: int(count)}}})
		}

	}

	return &dto, nil
}

func (r *reportRepository) GetReportViolationTrend(ctx context.Context, projectID uuid.UUID) (*models.ReportViolationTrendDTO, error) {
	trends := make([]models.ViolationTrend, 0)

	query := `
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
	`
	if err := r.db.WithContext(ctx).Raw(query, constants.CriticalDirectives, constants.HighDirectives, constants.MediumDirectives, projectID).Scan(&trends).Error; err != nil {
		return nil, err
	}

	dto := models.ReportViolationTrendDTO(trends)
	return &dto, nil
}

func (r *reportRepository) GetReportTopViolatedDirectives(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDirectivesDTO, error) {
	var dto models.ReportTopViolatedDirectivesDTO

	query := `
		WITH GroupedCounts AS (
			SELECT
					directive,
					COUNT(*) AS count
			FROM
					csp_reports
			WHERE
					project_id = ?
			GROUP BY
					directive
		),
		TotalCount AS (
			SELECT
				COUNT(*) AS total_violations
			FROM
				csp_reports
			WHERE
				project_id = ?
		)
		SELECT
				tc.total_violations,
				gc.directive,
				gc.count,
				(gc.count::numeric * 100.0 / tc.total_violations) AS percentage
		FROM
				GroupedCounts gc
		CROSS JOIN
				TotalCount tc
		ORDER BY
				gc.count DESC;
	`
	rows, err := r.db.WithContext(ctx).Raw(query, projectID, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var totalViolations int
		var directive string
		var count int
		var percentage float64
		if err := rows.Scan(&totalViolations, &directive, &count, &percentage); err != nil {
			return nil, err
		}
		dto.TotalViolations = totalViolations
		dto.Violations = append(dto.Violations, models.Violation{Directive: directive, Count: count, Percentage: percentage})
	}

	return &dto, nil
}

func (r *reportRepository) GetReportTopViolatedDocumentURLs(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolatedDocumentURLDTO, error) {
	var dto models.ReportTopViolatedDocumentURLDTO

	query := `
		SELECT
			document_url as url,
			COUNT(*) AS count
		FROM csp_reports
		WHERE project_id = ?
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY document_url
		ORDER BY count DESC
		LIMIT 10
	`
	if err := r.db.WithContext(ctx).Raw(query, projectID).Scan(&dto).Error; err != nil {
		return nil, err
	}

	return &dto, nil
}

func (r *reportRepository) GetReportSoftwareStats(ctx context.Context, projectID uuid.UUID) (*models.ReportSoftwareStatsDTO, error) {
	var dto models.ReportSoftwareStatsDTO

	query := `
		SELECT
			user_agent,
			COUNT(*) AS count
		FROM csp_reports
		WHERE project_id = ?
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY user_agent
		ORDER BY user_agent ASC
	`
	rows, err := r.db.WithContext(ctx).Raw(query, projectID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalCount int

	for rows.Next() {
		var userAgent string
		var count int
		if err := rows.Scan(&userAgent, &count); err != nil {
			return nil, err
		}

		totalCount += count

		browserOS := util.ParseUserAgent(userAgent)
		browserFound := false
		for i, item := range dto.Browser {
			if item.Name == browserOS.Browser {
				dto.Browser[i].Value += float64(count)
				browserFound = true
				break
			}
		}
		if !browserFound {
			dto.Browser = append(dto.Browser, models.StatItem{Name: browserOS.Browser, Value: float64(count)})
		}

		osFound := false
		for i, item := range dto.OS {
			if item.Name == browserOS.OS {
				dto.OS[i].Value += float64(count)
				osFound = true
				break
			}
		}
		if !osFound {
			dto.OS = append(dto.OS, models.StatItem{Name: browserOS.OS, Value: float64(count)})
		}
	}

	if totalCount > 0 {
		for i := range dto.Browser {
			val := (dto.Browser[i].Value / float64(totalCount)) * 100
			dto.Browser[i].Value = math.Round(val*100) / 100
		}
		for i := range dto.OS {
			val := (dto.OS[i].Value / float64(totalCount)) * 100
			dto.OS[i].Value = math.Round(val*100) / 100
		}
	}

	return &dto, nil
}

func (r *reportRepository) GetReportTopViolationSources(ctx context.Context, projectID uuid.UUID) (*models.ReportTopViolationSourcesDTO, error) {
	var dto models.ReportTopViolationSourcesDTO

	query := `
		WITH violations AS (
			SELECT
				blocked_url,
				directive,
				MAX(created_at) as last_seen,
				COUNT(*) AS count
			FROM csp_reports
			WHERE project_id = ?
			AND blocked_url IS NOT NULL AND blocked_url != ''
			AND created_at >= CURRENT_DATE - INTERVAL '30 days'
			GROUP BY blocked_url, directive
		),
		severity_map AS (
			SELECT
				directive,
				CASE
					WHEN directive IN ? THEN 'critical'
					WHEN directive IN ? THEN 'high'
					WHEN directive IN ? THEN 'medium'
					WHEN directive IN ? THEN 'low'
					ELSE 'unknown'
				END AS severity
			FROM (
				SELECT DISTINCT directive
				FROM csp_reports
				WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
			) d
		)
		SELECT
			v.blocked_url,
			SUM(v.count) AS count,
			s.severity,
			MAX(v.last_seen) AS last_seen
		FROM violations v
		JOIN severity_map s ON v.directive = s.directive
		GROUP BY v.blocked_url, s.severity
		ORDER BY count DESC
		LIMIT 10
	`
	if err := r.db.WithContext(ctx).Raw(query, projectID, constants.CriticalDirectives, constants.HighDirectives, constants.MediumDirectives, constants.LowDirectives).Scan(&dto).Error; err != nil {
		return nil, err
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
