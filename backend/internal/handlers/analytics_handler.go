package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/internal/util"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	reportService  services.ReportService
	projectService services.ProjectService
}

func NewAnalyticsHandler(r services.ReportService, p services.ProjectService) *AnalyticsHandler {
	return &AnalyticsHandler{
		reportService:  r,
		projectService: p,
	}
}

// GetReportSummaryStats godoc
// @Summary      Get report summary stats
// @Description  Get summary statistics for reports of a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Summary stats"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get report summary stats"
// @Router       /api/analytics/{projectID}/summary-stats [get]
func (h *AnalyticsHandler) GetReportSummaryStats(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportSummaryStats(ctx, projectID)
	if err != nil {
		log.Error("failed to get report summary stats", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report summary stats"))
		return
	}

	response.Success(c, http.StatusOK, "", data)
}

// GetReportGraphData godoc
// @Summary      Get report graph data
// @Description  Get graph data for reports of a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Graph data"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get report graph data"
// @Router       /api/analytics/{projectID}/graph-data [get]
func (h *AnalyticsHandler) GetReportGraphData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportGraphData(ctx, projectID)
	if err != nil {
		log.Error("failed to get report graph data", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report graph data"))
		return
	}

	response.Success(c, http.StatusOK, "", data)
}

// GetReportTopViolatedDirectives godoc
// @Summary      Get top violated directives
// @Description  Get top violated directives for a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Top violated directives"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get top violated directives"
// @Router       /api/analytics/{projectID}/violated-directives [get]
func (h *AnalyticsHandler) GetReportTopViolatedDirectives(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportTopViolatedDirectives(ctx, projectID)
	if err != nil {
		log.Error("failed to get report top violated directives", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report top violated directives"))
		return
	}

	response.Success(c, http.StatusOK, "", data)
}

// GetReportViolationTrend godoc
// @Summary      Get report violation trend
// @Description  Get violation trend for a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Violation trend"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get report violation trend"
// @Router       /api/analytics/{projectID}/violation-trend [get]
func (h *AnalyticsHandler) GetReportViolationTrend(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportViolationTrend(ctx, projectID)
	if err != nil {
		log.Error("failed to get report violation trend", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report violation trend"))
		return
	}

	response.Success(c, http.StatusOK, "", data)
}

// GetReportTopViolatedDocumentURLs godoc
// @Summary      Get top violated document URLs
// @Description  Get top violated document URLs for a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Top violated document URLs"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get top violated document URLs"
// @Router       /api/analytics/{projectID}/violated-document-urls [get]
func (h *AnalyticsHandler) GetReportTopViolatedDocumentURLs(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportTopViolatedDocumentURLs(ctx, projectID)
	if err != nil {
		log.Error("failed to get report top violated document uris", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report top violated document uris"))
		return
	}

	response.Success(c, http.StatusOK, "", data)
}

// GetReportSoftwareStats godoc
// @Summary      Get report software stats
// @Description  Get software statistics (browser/OS) for reports of a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Software stats"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get report software stats"
// @Router       /api/analytics/{projectID}/software-stats [get]
func (h *AnalyticsHandler) GetReportSoftwareStats(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportSoftwareStats(ctx, projectID)
	if err != nil {
		log.Error("failed to get report software stats", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report software stats"))
		return
	}
	response.Success(c, http.StatusOK, "", data)
}

// GetReportTopViolationSources godoc
// @Summary      Get top violation sources
// @Description  Get top violation sources for a project
// @Tags         analytics
// @Produce      json
// @Param        projectID   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Top violation sources"
// @Failure      400  {object}  map[string]interface{} "Invalid request parameters"
// @Failure      500  {object}  map[string]interface{} "Failed to get top violation sources"
// @Router       /api/analytics/{projectID}/violation-sources [get]
func (h *AnalyticsHandler) GetReportTopViolationSources(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)
	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.String("project_id", projectID.String()), zap.Error(err))
		c.Error(err)
		return
	}

	data, err := h.reportService.GetReportTopViolationSources(ctx, projectID)
	if err != nil {
		log.Error("failed to get report top violation sources", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to get report top violation sources"))
		return
	}
	response.Success(c, http.StatusOK, "", data)
}
