package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/internal/util"
	"go.uber.org/zap"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler
func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetProjectByID godoc
// @Summary      Get project by ID
// @Description  Get details of a project by its ID
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Project details"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /api/projects/{id} [get]
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)

	project, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to fetch project", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch project")
		return
	}

	if project == nil {
		log.Warn("project not found", zap.String("project_id", projectID.String()))
		response.ErrorResponse(c, http.StatusNotFound, "project not found")
		return
	}

	response.Success(c, http.StatusOK, "", project)
}

// CreateProject godoc
// @Summary      Create project
// @Description  Create a new project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        body  body      models.ProjectUpsertDTO  true  "project info"
// @Success      201   {object}  map[string]interface{} "project created successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      500   {object}  map[string]interface{} "failed to create project"
// @Router       /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	project := &models.ProjectUpsertDTO{}
	if err := c.BindJSON(project); err != nil {
		log.Warn("invalid request payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	err := h.projectService.CreateProject(ctx, project)
	if err != nil {
		log.Error("failed to create project", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to create project")
		return
	}

	response.Success(c, http.StatusCreated, "project created successfully", nil)
}

// UpdateProject godoc
// @Summary      Update project
// @Description  Update an existing project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        projectID    path      string                  true  "project ID"
// @Param        body  body      models.ProjectUpsertDTO  true  "project info"
// @Success      200   {object}  map[string]interface{} "project updated successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      500   {object}  map[string]interface{} "failed to update project"
// @Router       /api/projects/{projectID} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	id := uuid.MustParse(params.ProjectID)

	project := &models.ProjectUpsertDTO{}
	if err := c.BindJSON(project); err != nil {
		log.Warn("invalid request payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	updatedProject, err := h.projectService.UpdateProject(ctx, id, project)
	if err != nil {
		log.Error("failed to update project", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to update project")
		return
	}

	response.Success(c, http.StatusOK, "project updated successfully", updatedProject)
}

// DeleteProject godoc
// @Summary      Delete project
// @Description  Delete a project by its ID
// @Tags         projects
// @Produce      json
// @Param        projectID   path      string  true  "project ID"
// @Success      200  {object}  map[string]interface{} "project deleted successfully"
// @Failure      403  {object}  map[string]interface{} "failed to delete project"
// @Router       /api/projects/{projectID} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.ProjectParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	projectID := uuid.MustParse(params.ProjectID)

	if err := h.projectService.DeleteProject(ctx, projectID); err != nil {
		log.Error("failed to delete project", zap.Error(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "project deleted successfully", nil)
}

// ListProjects godoc
// @Summary      List projects
// @Description  Get a list of all projects
// @Tags         projects
// @Produce      json
// @Success      200  {object}  map[string]interface{} "List of projects"
// @Failure      403  {object}  map[string]interface{} "failed to list projects"
// @Router       /api/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)
	projects, err := h.projectService.ListProjects(ctx)
	if err != nil {
		log.Error("failed to list projects", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to list projects"))
		return
	}

	response.Success(c, http.StatusOK, "", projects)
}
