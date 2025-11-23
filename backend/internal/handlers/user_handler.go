package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/common"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
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

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get details of a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{} "User details"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	claims := common.GetUserClaims(c)

	user, err := h.userService.GetUserByID(id, claims)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.SuccessResponse(c, http.StatusOK, "", user)
}

// CreateUser godoc
// @Summary      Create user
// @Description  Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      models.UserCreateDTO  true  "User info"
// @Success      201   {object}  map[string]interface{} "User created successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to create user"
// @Router       /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.UserCreateDTO
	if err := c.BindJSON(&user); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	created, err := h.userService.RegisterUser(&user)
	if err != nil || !created {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	common.SuccessResponse(c, http.StatusCreated, "User created successfully", nil)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string              true  "User ID"
// @Param        body  body      models.UserUpdateDTO  true  "User info"
// @Success      200   {object}  map[string]interface{} "User updated successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to update user"
// @Router       /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.UserUpdateDTO
	if err := c.BindJSON(&user); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedUser, err := h.userService.UpdateUser(id, &user)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "User updated successfully", updatedUser)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{} "User deleted successfully"
// @Failure      500  {object}  map[string]interface{} "Failed to delete user"
// @Router       /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	claims := common.GetUserClaims(c)

	err := h.userService.DeleteUser(id, claims)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
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
	claims := common.GetUserClaims(c)

	users, err := h.userService.ListUsers(claims)

	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to list users")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "", users)
}

// ListUsersByTeamID godoc
// @Summary      List users by team ID
// @Description  Get a list of users by team ID
// @Tags         users
// @Produce      json
// @Param        teamID   path      string  true  "Team ID"
// @Success      200  {object}  map[string]interface{} "List of users by team"
// @Failure      500  {object}  map[string]interface{} "Failed to list users by team ID"
// @Router       /api/teams/{teamID}/users [get]
func (h *UserHandler) ListUsersByTeamID(c *gin.Context) {
	teamID := c.Param("teamID")

	claims := common.GetUserClaims(c)

	users, err := h.userService.ListUsersByTeamID(teamID, claims)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, "Failed to list users by team ID")
		return
	}

	common.SuccessResponse(c, http.StatusOK, "", users)
}
