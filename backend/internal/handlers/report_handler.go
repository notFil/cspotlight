package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/services"
	"go.uber.org/zap"
)

type ReportHandler struct {
	reportService services.ReportService
	reportQueue   chan reportJob
}

type reportJob struct {
	report    *models.CSPReportCreateDTO
	projectID string
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService services.ReportService) *ReportHandler {
	h := &ReportHandler{
		reportService: reportService,
		reportQueue:   make(chan reportJob, 1000),
	}
	go h.startWorker()
	return h
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	projectID := c.Param("projectID")

	report := models.CSPReportCreateDTO{}
	if err := c.BindJSON(&report); err != nil {
		log.Warn("invalid report payload", zap.Error(err))
		c.Error(err)
		return
	}

	report.SourceIP = c.ClientIP()

	select {
	case h.reportQueue <- reportJob{report: &report, projectID: projectID}:
	default:
		log.Warn("report queue full, dropping report", zap.String("project_id", projectID))
	}

	c.Status(http.StatusNoContent)
}

func (h *ReportHandler) ListReportsByProjectID(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	projectID := c.Param("projectID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	p := &pagination.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	log.Info("listing reports", zap.String("project_id", projectID))

	reports, meta, err := h.reportService.ListReportsByProjectID(ctx, projectID, p)
	if err != nil {
		log.Error("failed to list reports", zap.String("project_id", projectID), zap.Error(err))
		c.Error(err)
		return
	}
	response.SuccessPagedResponse(c, http.StatusOK, "reports fetched successfully", reports, meta)
}

func (h *ReportHandler) startWorker() {
	const batchSize = 100
	const batchTimeout = 5 * time.Second

	var batch []reportJob
	timer := time.NewTicker(batchTimeout)
	defer timer.Stop()

	for {
		select {
		case job := <-h.reportQueue:
			batch = append(batch, job)
			if len(batch) >= batchSize {
				h.processBatch(batch)
				batch = nil
				timer.Reset(batchTimeout)
			}
		case <-timer.C:
			if len(batch) > 0 {
				h.processBatch(batch)
				batch = nil
			}
		}
	}
}

func (h *ReportHandler) processBatch(batch []reportJob) {
	// Group reports by project ID to optimize DB calls

	grouped := make(map[string][]*models.CSPReportCreateDTO)
	for _, job := range batch {
		grouped[job.projectID] = append(grouped[job.projectID], job.report)
	}

	for projectID, reports := range grouped {
		if err := h.reportService.BatchCreateReports(context.Background(), reports, projectID); err != nil { // Passing nil context as c and claims are unavailable
			logger.Logger.Error("failed to create batch reports", zap.String("project_id", projectID), zap.Error(err))
		}
	}
}
