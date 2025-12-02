package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/logger"
	"github.com/notFil/cspotlight/pkg/response"
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
	id := c.Param("id")

	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	log.Info("fetching project", zap.String("project_id", id), zap.String("user_id", claims.Subject))

	project, err := h.projectService.GetProjectByID(id, *claims)
	if err != nil {
		log.Error("failed to fetch project", zap.Error(err))
		response.Error(c, err)
		return
	}
	if project == nil {
		log.Warn("project not found", zap.String("project_id", id))
		response.Error(c, err)
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
// @Param        body  body      models.ProjectUpsertDTO  true  "Project info"
// @Success      201   {object}  map[string]interface{} "Project created successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to create project"
// @Router       /api/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	project := &models.ProjectUpsertDTO{}
	if err := c.BindJSON(project); err != nil {
		log.Warn("invalid request payload", zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("creating project", zap.String("user_id", claims.Subject), zap.String("project_name", project.Name))

	created, err := h.projectService.CreateProject(project, *claims)
	if err != nil || !created {
		log.Error("failed to create project", zap.Error(err))
		response.Error(c, err)
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
// @Param        id    path      string                  true  "Project ID"
// @Param        body  body      models.ProjectUpsertDTO  true  "Project info"
// @Success      200   {object}  map[string]interface{} "Project updated successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to update project"
// @Router       /api/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")

	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	project := &models.ProjectUpsertDTO{}
	if err := c.BindJSON(project); err != nil {
		log.Warn("invalid request payload", zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("updating project", zap.String("project_id", id), zap.String("user_id", claims.Subject))

	updatedProject, err := h.projectService.UpdateProject(id, project, *claims)
	if err != nil {
		log.Error("failed to update project", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "", updatedProject)
}

// DeleteProject godoc
// @Summary      Delete project
// @Description  Delete a project by its ID
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Project deleted successfully"
// @Failure      403  {object}  map[string]interface{} "failed to delete project"
// @Router       /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	log.Info("deleting project", zap.String("project_id", id), zap.String("user_id", claims.Subject))

	if err := h.projectService.DeleteProject(id, *claims); err != nil {
		log.Error("failed to delete project", zap.Error(err))
		response.Error(c, err)
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
	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	log.Info("listing projects", zap.String("user_id", claims.Subject))

	projects, err := h.projectService.ListProjects(*claims)
	if err != nil {
		log.Error("failed to list projects", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "", projects)
}
