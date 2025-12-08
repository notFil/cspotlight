package handlers

import (
	"net/http"

	"github.com/notFil/cspotlight/internal/auth"
	"github.com/notFil/cspotlight/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(s services.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Get details of the current user
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]interface{} "User details"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /api/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	claims := auth.GetUserClaims(ctx)
	log.Info("fetching current user", zap.String("user_id", claims.Subject))

	user, err := h.userService.GetUserByID(ctx, claims.Subject)
	if err != nil {
		log.Error("failed to fetch current user", zap.String("user_id", claims.Subject), zap.Error(err))
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "", user)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get details of a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{} "user details"
// @Failure      500  {object}  map[string]interface{} "internal server error"
// @Router       /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	id := c.Param("id")

	log.Info("fetching user", zap.String("user_id", id))

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		log.Error("failed to fetch user", zap.String("user_id", id), zap.Error(err))
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "", user)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        body  body      models.UserUpdateDTO  true  "User info"
// @Success      200   {object}  map[string]interface{} "user updated successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      500   {object}  map[string]interface{} "failed to update user"
// @Router       /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	id := c.Param("id")
	var user models.UserUpdateDTO
	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid user update payload", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("updating user", zap.String("user_id", id))

	updatedUser, err := h.userService.UpdateUser(ctx, id, &user)
	if err != nil {
		log.Error("failed to update user", zap.String("user_id", id), zap.Error(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "user updated successfully", updatedUser)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{} "user deleted successfully"
// @Failure      500  {object}  map[string]interface{} "failed to delete user"
// @Router       /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	id := c.Param("id")

	log.Info("deleting user", zap.String("user_id", id))

	err := h.userService.DeleteUser(ctx, id)
	if err != nil {
		log.Error("failed to delete user", zap.String("user_id", id), zap.Error(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "user deleted successfully", nil)
}

// ListUsers godoc
// @Summary      List users
// @Description  Get a list of all users
// @Tags         users
// @Produce      json
// @Success      200  {object}  map[string]interface{} "List of users"
// @Failure      500  {object}  map[string]interface{} "Failed to list users"
// @Router       /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	log.Info("listing users")
	users, err := h.userService.ListUsers(ctx)

	if err != nil {
		log.Error("failed to list users", zap.Error(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "", users)
}

// ListUsersByTeamID godoc
// @Summary      List users by team ID
// @Description  Get a list of users by team ID
// @Tags         users
// @Produce      json
// @Param        teamID   path      string  true  "Team ID"
// @Success      200  {object}  map[string]interface{} "list of users by team"
// @Failure      500  {object}  map[string]interface{} "failed to list users by team ID"
// @Router       /api/teams/{teamID}/users [get]
func (h *UserHandler) ListUsersByTeamID(c *gin.Context) {
	teamID := c.Param("teamID")

	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	log.Info("listing users by team", zap.String("team_id", teamID))
	users, err := h.userService.ListUsersByTeamID(ctx, teamID)
	if err != nil {
		log.Error("failed to list users by team", zap.String("team_id", teamID), zap.Error(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "", users)
}

// SetDefaultProject godoc
// @Summary      Set default project
// @Description  Sets the default project for a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        body  body      models.SetDefaultProjectRequest  true  "Set default project request"
// @Success      200   {object}  map[string]interface{} "default project set successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      500   {object}  map[string]interface{} "failed to set default project"
// @Router       /api/users/{id}/project [patch]
func (h *UserHandler) SetDefaultProject(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	log := logger.FromContext(ctx)
	var req struct {
		ProjectID string `json:"projectID"`
	}
	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid set default project payload", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("setting default project", zap.String("project_id", req.ProjectID), zap.String("user_id", id))
	updatedUser, err := h.userService.SetDefaultProject(ctx, id, req.ProjectID)
	if err != nil {
		log.Error("failed to set default project", zap.Error(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "default project set successfully", updatedUser)

}

// ChangePassword godoc
// @Summary      Change password
// @Description  Changes the password of a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        body  body      models.ChangePasswordRequest  true  "Change password request"
// @Success      200   {object}  map[string]interface{} "password changed successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      400   {object}  map[string]interface{} "failed to change password"
// @Router       /api/users/{id}/password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	id := c.Param("id")

	var req models.ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid change password payload", zap.Error(err))
		c.Error(err)
		return
	}

	err := h.userService.ChangePassword(ctx, id, &req)
	if err != nil {
		log.Error("failed to change password", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("user changed password", zap.String("user_id", id))

	response.Success(c, http.StatusOK, "password changed successfully", nil)
}

// ChangeImage godoc
// @Summary      Change image
// @Description  Changes the image of a user
// @Tags         users
// @Accept       multipart
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        image  body      models.ChangeImageRequest  true  "Change image request"
// @Success      200   {object}  map[string]interface{} "image changed successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      400   {object}  map[string]interface{} "failed to change image"
// @Router       /api/users/{id}/image [patch]
func (h *UserHandler) ChangeImage(c *gin.Context) {
	id := c.Param("id")

	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	file, err := c.FormFile("image")
	if err != nil {
		log.Warn("invalid change image payload", zap.Error(err))
		c.Error(err)
		return
	}

	err = h.userService.ChangeImage(ctx, id, file)
	if err != nil {
		log.Error("failed to change image", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("user changed image", zap.String("user_id", id))

	response.Success(c, http.StatusOK, "image changed successfully", nil)
}
