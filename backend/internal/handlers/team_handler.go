package handlers

import (
	"net/http"

	"github.com/google/uuid"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
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
// @Param        teamID   path      string  true  "Team ID"
// @Success      200  {object}  map[string]interface{} "Team details"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /api/teams/{teamID} [get]
func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.TeamParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	teamID := uuid.MustParse(params.TeamID)

	log.Info("fetching team", zap.String("team_id", teamID.String()))

	user, err := h.teamService.GetTeamByID(ctx, teamID)
	if err != nil {
		log.Error("failed to fetch team", zap.String("team_id", teamID.String()), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch team")
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
// @Param        body  body      models.TeamUpsertDTO  true  "team info"
// @Success      201   {object}  map[string]interface{} "team created successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      500   {object}  map[string]interface{} "failed to create team"
// @Router       /api/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	team := &models.TeamUpsertDTO{}
	if err := c.BindJSON(team); err != nil {
		log.Warn("invalid team payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	log.Info("creating team")

	if err := h.teamService.CreateTeam(ctx, team); err != nil {
		log.Error("failed to create team", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to create team")
		return
	}

	response.Success(c, http.StatusCreated, "team created successfully", nil)
}

// UpdateTeam godoc
// @Summary      Update team
// @Description  Update an existing team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        teamID    path      string              true  "Team ID"
// @Param        body  body      models.TeamUpsertDTO  true  "Team info"
// @Success      200   {object}  map[string]interface{} "Team updated successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to update team"
// @Router       /api/teams/{teamID} [put]
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.TeamParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	teamID := uuid.MustParse(params.TeamID)

	team := &models.TeamUpsertDTO{}
	if err := c.BindJSON(team); err != nil {
		log.Warn("invalid team update payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info("updating team", zap.String("team_id", teamID.String()))

	updatedTeam, err := h.teamService.UpdateTeam(ctx, teamID, team)
	if err != nil {
		log.Error("failed to update team", zap.String("team_id", teamID.String()), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update team")
		return
	}

	response.Success(c, http.StatusOK, "", updatedTeam)
}

// DeleteTeam godoc
// @Summary      Delete team
// @Description  Delete a team by its ID
// @Tags         teams
// @Produce      json
// @Param        teamID   path      string  true  "Team ID"
// @Success      200  {object}  map[string]interface{} "Team deleted successfully"
// @Failure      500  {object}  map[string]interface{} "Failed to delete team"
// @Router       /api/teams/{teamID} [delete]
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.TeamParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	teamID := uuid.MustParse(params.TeamID)

	log.Info("deleting team", zap.String("team_id", teamID.String()))

	if err := h.teamService.DeleteTeam(ctx, teamID); err != nil {
		log.Error("failed to delete team", zap.String("team_id", teamID.String()), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete team")
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
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	log.Info("listing teams")

	teams, err := h.teamService.ListTeams(ctx)
	if err != nil {
		log.Error("failed to list teams", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to list teams")
		return
	}

	response.Success(c, http.StatusOK, "", teams)
}
