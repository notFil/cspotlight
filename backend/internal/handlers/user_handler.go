package handlers

import (
	"net/http"

	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/errs"
	"github.com/notFil/cspotlight/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/logger"
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
	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)
	log.Info("fetching current user", zap.String("user_id", claims.Subject))

	user, err := h.userService.GetUserByID(claims.Subject, *claims)
	if err != nil {
		log.Error("failed to fetch current user", zap.String("user_id", claims.Subject), zap.Error(err))
		response.Error(c, err)
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
	id := c.Param("id")

	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	log.Info("fetching user", zap.String("target_user_id", id), zap.String("user_id", claims.Subject))

	user, err := h.userService.GetUserByID(id, *claims)
	if err != nil {
		log.Error("failed to fetch user", zap.String("target_user_id", id), zap.Error(err))
		response.Error(c, err)
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
	id := c.Param("id")
	log := logger.FromContext(c)
	var user models.UserUpdateDTO
	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid user update payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Info("updating user", zap.String("target_user_id", id), zap.String("user_id", id))

	updatedUser, err := h.userService.UpdateUser(id, &user)
	if err != nil {
		log.Error("failed to update user", zap.String("target_user_id", id), zap.Error(err))
		response.Error(c, err)
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
	id := c.Param("id")

	log := logger.FromContext(c)

	log.Info("deleting user", zap.String("target_user_id", id), zap.String("user_id", id))

	err := h.userService.DeleteUser(id)
	if err != nil {
		log.Error("failed to delete user", zap.String("target_user_id", id), zap.Error(err))
		response.Error(c, err)
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
	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	log.Info("listing users", zap.String("user_id", claims.Subject))

	users, err := h.userService.GetUsers(*claims)

	if err != nil {
		log.Error("failed to list users", zap.Error(err))
		response.Error(c, err)
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

	claims := auth.GetUserClaims(c)
	log := logger.FromContext(c)

	log.Info("listing users by team", zap.String("team_id", teamID), zap.String("user_id", claims.Subject))

	users, err := h.userService.ListUsersByTeamID(teamID, *claims)
	if err != nil {
		log.Error("failed to list users by team", zap.String("team_id", teamID), zap.Error(err))
		response.Error(c, err)
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

	log := logger.FromContext(c)
	var req struct {
		ProjectID string `json:"projectID"`
	}
	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid set default project payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	claims := auth.GetUserClaims(c)

	log.Info("setting default project", zap.String("project_id", req.ProjectID), zap.String("user_id", claims.Subject))

	updatedUser, err := h.userService.SetDefaultProject(id, req.ProjectID, *claims)
	if err != nil {
		log.Error("failed to set default project", zap.Error(err))
		response.Error(c, err)
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
	id := c.Param("id")
	claims := auth.GetUserClaims(c)

	log := logger.FromContext(c)
	var req models.ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid change password payload", zap.Error(err))
		response.Error(c, err)
		return
	}
	if req.NewPassword != req.ConfirmNewPassword {
		log.Warn("new password and confirm new password do not match")
		response.Error(c, errs.ErrNoMatchOnNewPasswords)
		return
	}

	err := h.userService.ChangePassword(id, &req, *claims)
	if err != nil {
		log.Error("failed to change password", zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("user changed password", zap.String("user_id", claims.Subject))

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
	claims := auth.GetUserClaims(c)

	log := logger.FromContext(c)

	file, err := c.FormFile("image")
	if err != nil {
		log.Warn("invalid change image payload", zap.Error(err))
		response.Error(c, errs.ErrMissingImage)
		return
	}

	err = h.userService.ChangeImage(id, file, *claims)
	if err != nil {
		log.Error("failed to change image", zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("user changed image", zap.String("user_id", claims.Subject))

	response.Success(c, http.StatusOK, "image changed successfully", nil)
}
