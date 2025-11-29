package handlers

import (
	"net/http"

	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
)

type AuthHandler struct {
	userService services.UserService
	jwtConfig   config.JWTConfig
}

func NewAuthHandler(userService services.UserService, jwtConfig config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtConfig:   jwtConfig,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.AuthRequest  true  "Login credentials"
// @Success      200   {object}  map[string]interface{} "Login successful"
// @Failure      401   {object}  map[string]interface{} "Authentication failed"
// @Failure      500   {object}  map[string]interface{} "Internal server error"
// @Router       /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var authRequest models.AuthRequest
	if err := c.BindJSON(&authRequest); err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Authentication failed")
		return
	}

	user, err := h.userService.AuthenticateUser(&authRequest)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed")
		return
	}

	tokens, err := auth.GenerateJWT(user, h.jwtConfig)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Authentication failed")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Login successful", tokens)
}

// Register godoc
// @Summary      Register user
// @Description  Registers a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.UserCreateDTO  true  "User registration info"
// @Success      201   {object}  map[string]interface{} "User registered successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to register user"
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.UserCreateDTO
	if err := c.BindJSON(&user); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err := h.userService.RegisterUser(&user)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	response.SuccessResponse(c, http.StatusCreated, "User registered successfully", nil)
}

// Refresh token godoc
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, err := auth.ValidateRefreshToken(req.RefreshToken, h.jwtConfig.RefreshSecretKey)
	if err != nil {
		response.ErrorResponse(c, http.StatusUnauthorized, "Failed to refresh token")
		return
	}

	authClaims := auth.Claims{
		UserID: claims.Subject,
	}
	user, err := h.userService.GetUserByID(claims.Subject, authClaims)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token")
		return
	}

	tokens, err := auth.GenerateJWT(user, h.jwtConfig)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Token refreshed", tokens)
}

// Sign out godoc
func (h *AuthHandler) SignOut(c *gin.Context) {

}
