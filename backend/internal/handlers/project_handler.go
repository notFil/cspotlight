package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/common"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
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

	claims := common.GetUserClaims(c)

	user, err := h.projectService.GetProjectByID(id, claims)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SuccessResponse(c, http.StatusOK, "", user)
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
	claims := common.GetUserClaims(c)

	project := &models.ProjectUpsertDTO{}
	if err := c.BindJSON(project); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	created, err := h.projectService.CreateProject(project, claims)
	if err != nil || !created {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to create project")
		return
	}

	common.SuccessResponse(c, http.StatusCreated, "Project created successfully", nil)
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

	claims := common.GetUserClaims(c)

	project := &models.ProjectUpsertDTO{}
	if err := c.BindJSON(project); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedProject, err := h.projectService.UpdateProject(id, project, claims)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to update project")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "", updatedProject)
}

// DeleteProject godoc
// @Summary      Delete project
// @Description  Delete a project by its ID
// @Tags         projects
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  map[string]interface{} "Project deleted successfully"
// @Failure      500  {object}  map[string]interface{} "Failed to delete project"
// @Router       /api/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	claims := common.GetUserClaims(c)

	if err := h.projectService.DeleteProject(id, claims); err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "Project deleted successfully", nil)
}

// ListProjects godoc
// @Summary      List projects
// @Description  Get a list of all projects
// @Tags         projects
// @Produce      json
// @Success      200  {object}  map[string]interface{} "List of projects"
// @Failure      500  {object}  map[string]interface{} "Failed to list projects"
// @Router       /api/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	claims := common.GetUserClaims(c)

	projects, err := h.projectService.ListProjects(claims)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to list projects")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "", projects)
}
