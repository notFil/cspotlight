package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService    services.UserService
	projectService services.ProjectService
	staticPath     string
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(u services.UserService, p services.ProjectService, staticPath string) *UserHandler {
	return &UserHandler{
		userService:    u,
		projectService: p,
		staticPath:     staticPath,
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

	userContext := auth.GetUserContext(ctx)
	log.Info("fetching current user", zap.String("user_id", userContext.ID.String()))

	user, err := h.userService.GetUserByID(ctx, userContext.ID)
	if err != nil {
		log.Error("failed to fetch current user", zap.String("user_id", userContext.ID.String()), zap.Error(err))
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
// @Param        userID   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{} "user details"
// @Failure      500  {object}  map[string]interface{} "internal server error"
// @Router       /api/users/{userID} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.UserParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(err)
		return
	}

	userID := uuid.MustParse(params.UserID)

	log.Info("fetching user", zap.String("user_id", userID.String()))

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("failed to fetch user", zap.String("user_id", userID.String()), zap.Error(err))
		if err.Error() == "user not found" {
			c.Error(apperrors.New(http.StatusNotFound, "user not found"))
			return
		}
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
// @Router       /api/users/{userID} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.UserParams

	var user models.UserUpdateDTO

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid user update payload", zap.Error(err))
		c.Error(err)
		return
	}

	userID := uuid.MustParse(params.UserID)

	log.Info("updating user", zap.String("user_id", userID.String()))

	updatedUser, err := h.userService.UpdateUser(ctx, userID, &user)
	if err != nil {
		log.Error("failed to update user", zap.String("user_id", userID.String()), zap.Error(err))
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
// @Router       /api/users/{userID} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.UserParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	userID := uuid.MustParse(params.UserID)

	log.Info("deleting user", zap.String("user_id", userID.String()))

	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		log.Error("failed to delete user", zap.String("user_id", userID.String()), zap.Error(err))
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
	var params util.TeamParams

	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	teamID := uuid.MustParse(params.TeamID)

	log.Info("listing users by team", zap.String("team_id", teamID.String()))
	users, err := h.userService.ListUsersByTeamID(ctx, teamID)
	if err != nil {
		log.Error("failed to list users by team", zap.String("team_id", teamID.String()), zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to list users by team"))
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
// @Param        userID    path      string              true  "User ID"
// @Param        body  body      models.SetDefaultProjectRequest  true  "Set default project request"
// @Success      200   {object}  map[string]interface{} "default project set successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      500   {object}  map[string]interface{} "failed to set default project"
// @Router       /api/users/{userID}/project [patch]
func (h *UserHandler) SetDefaultProject(c *gin.Context) {
	var params util.UserParams
	var req models.SetDefaultProjectRequest

	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	userID := uuid.MustParse(params.UserID)

	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid set default project payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request payload"))
		return
	}

	projectID := uuid.MustParse(req.ProjectID)
	_, err := h.projectService.GetProjectByID(ctx, projectID)
	if err != nil {
		log.Error("failed to get project", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("setting default project", zap.String("project_id", req.ProjectID), zap.String("user_id", userID.String()))
	updatedUser, err := h.userService.SetDefaultProject(ctx, userID, projectID)
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
// @Param        userID    path      string              true  "User ID"
// @Param        body  body      models.ChangePasswordRequest  true  "Change password request"
// @Success      200   {object}  map[string]interface{} "password changed successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      400   {object}  map[string]interface{} "failed to change password"
// @Router       /api/users/{userID}/password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	var params util.UserParams

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	userID := uuid.MustParse(params.UserID)

	var req models.ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid change password payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request payload"))
		return
	}

	if err := h.userService.ChangePassword(ctx, userID, &req); err != nil {
		log.Error("failed to change password", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("user changed password", zap.String("user_id", userID.String()))
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: -1})
	session.Clear()
	if err := session.Save(); err != nil {
		log.Error("failed to save session", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to sign out"))
		return
	}

	response.Success(c, http.StatusOK, "password changed successfully", nil)
}

// ChangeImage godoc
// @Summary      Change image
// @Description  Changes the image of a user
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        userID    path      string              true  "User ID"
// @Param        image     formData  file                        true  "Image file"
// @Success      200   {object}  map[string]interface{} "image changed successfully"
// @Failure      400   {object}  map[string]interface{} "invalid request payload"
// @Failure      400   {object}  map[string]interface{} "failed to change image"
// @Router       /api/users/{userID}/image [patch]
func (h *UserHandler) ChangeImage(c *gin.Context) {
	var params util.UserParams

	ctx := c.Request.Context()

	log := logger.FromContext(ctx)

	if err := c.ShouldBindUri(&params); err != nil {
		log.Error("failed to bind uri", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	userID := uuid.MustParse(params.UserID)

	file, err := c.FormFile("image")
	if err != nil {
		log.Warn("invalid change image payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid request payload"))
		return
	}

	uploadDir := fmt.Sprintf("%s/images/users", h.staticPath)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Error("failed to create upload directory", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to process image"))
		return
	}

	filename := fmt.Sprintf("%s_%d.png", uuid.New().String(), time.Now().Unix())
	dst := filepath.Join(uploadDir, filename)

	if err := util.ValidateImage(file); err != nil {
		log.Error("failed to validate image", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid image"))
		return
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		log.Error("failed to save uploaded file", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to save uploaded file"))
		return
	}

	err = h.userService.ChangeImage(ctx, userID, dst)
	if err != nil {
		log.Error("failed to change image", zap.Error(err))
		c.Error(err)
		return
	}

	log.Info("user changed image", zap.String("user_id", userID.String()))

	response.Success(c, http.StatusOK, "image changed successfully", nil)
}
