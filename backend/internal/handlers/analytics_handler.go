package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/logger"
	"github.com/notFil/cspotlight/pkg/response"
)

type AnalyticsHandler struct {
	reportService services.ReportService
}

func NewAnalyticsHandler(s services.ReportService) *AnalyticsHandler {
	return &AnalyticsHandler{
		reportService: s,
	}
}

func (h *AnalyticsHandler) GetReportGraphData(c *gin.Context) {
	claims := auth.GetUserClaims(c)

	log := logger.FromContext(c)

	projectID := c.Param("projectID")

	data, err := h.reportService.GetReportGraphData(projectID, *claims)
	if err != nil {
		log.Error("failed to get report graph data")
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "", data)
}
