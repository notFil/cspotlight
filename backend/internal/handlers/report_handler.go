package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	apperrors "cspotlight/internal/errors"
	"cspotlight/internal/logger"
	"cspotlight/internal/models"
	"cspotlight/internal/pagination"
	"cspotlight/internal/response"
	"cspotlight/internal/services"
	"cspotlight/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReportHandler struct {
	reportService  services.ReportService
	projectService services.ProjectService
	reportQueue    chan reportJob
}

type reportJob struct {
	report    *models.CSPReportCreateDTO
	projectID uuid.UUID
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService services.ReportService, projectService services.ProjectService) *ReportHandler {
	h := &ReportHandler{
		reportService:  reportService,
		projectService: projectService,
		reportQueue:    make(chan reportJob, 1000),
	}
	go h.startWorker()
	return h
}

// CreateReport godoc
// @Summary      Create report
// @Description  Create a new CSP report
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Param        body  body      models.CSPReportCreateDTO  true  "Report info"
// @Success      204   "No Content"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to create report"
// @Router       /api/reports/{projectID}/endpoint [post]
func (h *ReportHandler) CreateReport(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)

	report := models.CSPReportCreateDTO{}
	if err := c.BindJSON(&report); err != nil {
		log.Warn("invalid report payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	report.SourceIP = c.ClientIP()

	select {
	case h.reportQueue <- reportJob{report: &report, projectID: projectID}:
	default:
		log.Warn("report queue full, dropping report", zap.String("project_id", projectID.String()))
	}

	c.Status(http.StatusNoContent)
}

// ListReportsByProjectID godoc
// @Summary      List reports by project ID
// @Description  Get a list of reports for a project
// @Tags         reports
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Param        page        query     int     false "Page number"
// @Param        page_size   query     int     false "Page size"
// @Success      200  {object}  map[string]interface{} "List of reports"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to list reports"
// @Router       /api/reports/{projectID} [get]
func (h *ReportHandler) ListReportsByProjectID(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	p := &pagination.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	projectID := uuid.MustParse(params.ProjectID)

	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("listing reports", zap.String("project_id", projectID.String()))

	reports, meta, err := h.reportService.ListReportsByProjectID(ctx, projectID, p)
	if err != nil {
		log.Error("failed to list reports", zap.String("project_id", projectID.String()), zap.Error(err))
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

	grouped := make(map[uuid.UUID][]*models.CSPReportCreateDTO)
	for _, job := range batch {
		grouped[job.projectID] = append(grouped[job.projectID], job.report)
	}

	for projectID, reports := range grouped {
		if err := h.reportService.BatchCreateReports(context.Background(), reports, projectID); err != nil { // Passing nil context as c and claims are unavailable
			logger.Logger.Error("failed to create batch reports", zap.String("project_id", projectID.String()), zap.Error(err))
		}
	}
}
