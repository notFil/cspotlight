package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/pagination"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/response"
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
	projectID := c.Param("projectID")

	report := models.CSPReportCreateDTO{}
	if err := c.BindJSON(&report); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := h.reportService.CreateReport(&report, projectID); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create report")
		return
	}
	response.SuccessResponse(c, http.StatusCreated, "Report created successfully", nil)
}

func (h *ReportHandler) ListReportsByProjectID(c *gin.Context) {
	claims := auth.GetUserClaims(c)

	projectID := c.Param("projectID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	p := &pagination.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	reports, meta, err := h.reportService.ListReportsByProjectID(projectID, p, *claims)
	if err != nil {
		response.ErrorResponse(c, http.StatusNotFound, "Failed to list reports")
	}
	response.SuccessPagedResponse(c, http.StatusOK, "Reports fetched successfully", reports, meta)
}
