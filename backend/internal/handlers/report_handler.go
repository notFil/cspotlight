package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/logger"
	"github.com/notFil/cspotlight/pkg/response"
	"go.uber.org/zap"
)

type ReportHandler struct {
	reportService services.ReportService
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	log := logger.FromContext(c)
	projectID := c.Param("projectID")

	report := models.CSPReportCreateDTO{}
	if err := c.BindJSON(&report); err != nil {
		log.Warn("invalid report payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info("creating report", zap.String("project_id", projectID))

	if err := h.reportService.CreateReport(&report, projectID); err != nil {
		log.Error("failed to create report", zap.String("project_id", projectID), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create report")
		return
	}
	response.SuccessResponse(c, http.StatusCreated, "Report created successfully", nil)
}

func (h *ReportHandler) ListReportsByProjectID(c *gin.Context) {
	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	projectID := c.Param("projectID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	p := &pagination.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	log.Info("listing reports", zap.String("project_id", projectID), zap.String("user_id", claims.UserID))

	reports, meta, err := h.reportService.ListReportsByProjectID(projectID, p, *claims)
	if err != nil {
		log.Error("failed to list reports", zap.String("project_id", projectID), zap.Error(err))
		response.ErrorResponse(c, http.StatusNotFound, "Failed to list reports")
	}
	response.SuccessPagedResponse(c, http.StatusOK, "Reports fetched successfully", reports, meta)
}
