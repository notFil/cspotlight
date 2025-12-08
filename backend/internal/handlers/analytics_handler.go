package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/services"
	"go.uber.org/zap"
)

type AnalyticsHandler struct {
	reportService services.ReportService
}

func NewAnalyticsHandler(s services.ReportService) *AnalyticsHandler {
	return &AnalyticsHandler{
		reportService: s,
	}
}

func (h *AnalyticsHandler) GetReportSummaryStats(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	projectID := c.Param("projectID")

	data, err := h.reportService.GetReportSummaryStats(ctx, projectID)
	if err != nil {
		log.Error("failed to get report summary stats", zap.Error(err))
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "", data)
}

func (h *AnalyticsHandler) GetReportGraphData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	projectID := c.Param("projectID")

	data, err := h.reportService.GetReportGraphData(ctx, projectID)
	if err != nil {
		log.Error("failed to get report graph data", zap.Error(err))
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "", data)
}

func (h *AnalyticsHandler) GetReportViolationTrend(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	projectID := c.Param("projectID")

	data, err := h.reportService.GetReportViolationTrend(ctx, projectID)
	if err != nil {
		log.Error("failed to get report violation trend", zap.Error(err))
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "", data)
}

func (h *AnalyticsHandler) GetReportSoftwareStats(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	projectID := c.Param("projectID")

	data, err := h.reportService.GetReportSoftwareStats(ctx, projectID)
	if err != nil {
		log.Error("failed to get report software stats", zap.Error(err))
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "", data)
}
