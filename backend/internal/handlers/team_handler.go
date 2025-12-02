package handlers

import (
	"net/http"

	"github.com/notFil/cspotlight/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/logger"
	"go.uber.org/zap"
)

type TeamHandler struct {
	teamService services.TeamService
}

// NewTeamHandler creates a new TeamHandler
func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// GetTeamByID godoc
// @Summary      Get team by ID
// @Description  Get details of a team by its ID
// @Tags         teams
// @Produce      json
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  map[string]interface{} "Team details"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /api/teams/{id} [get]
func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	log := logger.FromContext(c)
	id := c.Param("id")

	log.Info("fetching team", zap.String("team_id", id))

	user, err := h.teamService.GetTeamByID(id)
	if err != nil {
		log.Error("failed to fetch team", zap.String("team_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "", user)
}

// CreateTeam godoc
// @Summary      Create team
// @Description  Create a new team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        body  body      models.TeamUpsertDTO  true  "Team info"
// @Success      201   {object}  map[string]interface{} "Team created successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to create team"
// @Router       /api/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	log := logger.FromContext(c)
	team := &models.TeamUpsertDTO{}
	if err := c.BindJSON(team); err != nil {
		log.Warn("invalid team payload", zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("creating team", zap.String("team_name", team.Name))

	if err := h.teamService.CreateTeam(team); err != nil {
		log.Error("failed to create team", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Team created successfully", nil)
}

// UpdateTeam godoc
// @Summary      Update team
// @Description  Update an existing team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "Team ID"
// @Param        body  body      models.TeamUpsertDTO  true  "Team info"
// @Success      200   {object}  map[string]interface{} "Team updated successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to update team"
// @Router       /api/teams/{id} [put]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	log := logger.FromContext(c)
	id := c.Param("id")
	team := &models.TeamUpsertDTO{}
	if err := c.BindJSON(team); err != nil {
		log.Warn("invalid team update payload", zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("updating team", zap.String("team_id", id))

	updatedTeam, err := h.teamService.UpdateTeam(id, team)
	if err != nil {
		log.Error("failed to update team", zap.String("team_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "", updatedTeam)
}

// DeleteTeam godoc
// @Summary      Delete team
// @Description  Delete a team by its ID
// @Tags         teams
// @Produce      json
// @Param        id   path      string  true  "Team ID"
// @Success      200  {object}  map[string]interface{} "Team deleted successfully"
// @Failure      500  {object}  map[string]interface{} "Failed to delete team"
// @Router       /api/teams/{id} [delete]
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	log := logger.FromContext(c)
	id := c.Param("id")

	log.Info("deleting team", zap.String("team_id", id))

	if err := h.teamService.DeleteTeam(id); err != nil {
		log.Error("failed to delete team", zap.String("team_id", id), zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Team deleted successfully", nil)
}

// ListTeams godoc
// @Summary      List teams
// @Description  Get a list of all teams
// @Tags         teams
// @Produce      json
// @Success      200  {object}  map[string]interface{} "List of teams"
// @Failure      500  {object}  map[string]interface{} "Failed to list teams"
// @Router       /api/teams [get]
func (h *TeamHandler) ListTeams(c *gin.Context) {
	log := logger.FromContext(c)
	log.Info("listing teams")

	teams, err := h.teamService.ListTeams()
	if err != nil {
		log.Error("failed to list teams", zap.Error(err))
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "", teams)
}
