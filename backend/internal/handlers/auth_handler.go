package handlers

import (
	"net/http"

	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/services"
	"github.com/notFil/cspotlight/pkg/logger"
	"go.uber.org/zap"
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
	log := logger.FromContext(c)
	var authRequest models.AuthRequest
	if err := c.BindJSON(&authRequest); err != nil {
		log.Warn("invalid login payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "Authentication failed")
		return
	}

	log.Info("attempting login", zap.String("username", authRequest.Username))

	user, err := h.userService.AuthenticateUser(&authRequest)
	if err != nil {
		log.Warn("authentication failed", zap.String("username", authRequest.Username), zap.Error(err))
		response.ErrorResponse(c, http.StatusUnauthorized, "Authentication failed")
		return
	}

	tokens, err := auth.GenerateJWT(user, h.jwtConfig)
	if err != nil {
		log.Error("failed to generate jwt", zap.String("user_id", user.ID), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "authentication failed")
		return
	}

	log.Info("login successful", zap.String("user_id", user.ID))
	response.SuccessResponse(c, http.StatusOK, "login successful", tokens)
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
	log := logger.FromContext(c)
	var user models.UserCreateDTO
	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid register payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	log.Info("registering user", zap.String("email", user.Email))

	_, err := h.userService.RegisterUser(&user)
	if err != nil {
		log.Error("failed to register user", zap.String("email", user.Email), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to register user")
		return
	}

	log.Info("user registered successfully", zap.String("email", user.Email))
	response.SuccessResponse(c, http.StatusCreated, "user registered successfully", nil)
}

// Refresh token godoc
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	log := logger.FromContext(c)
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid refresh token payload", zap.Error(err))
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	claims, err := auth.ValidateRefreshToken(req.RefreshToken, h.jwtConfig.RefreshSecretKey)
	if err != nil {
		log.Warn("invalid refresh token", zap.Error(err))
		response.ErrorResponse(c, http.StatusUnauthorized, "failed to refresh token")
		return
	}

	authClaims := auth.Claims{
		UserID: claims.Subject,
	}
	user, err := h.userService.GetUserByID(claims.Subject, authClaims)
	if err != nil {
		log.Error("failed to get user for refresh token", zap.String("user_id", claims.Subject), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	tokens, err := auth.GenerateJWT(user, h.jwtConfig)
	if err != nil {
		log.Error("failed to generate jwt for refresh", zap.String("user_id", user.ID), zap.Error(err))
		response.ErrorResponse(c, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	log.Info("token refreshed", zap.String("user_id", user.ID))
	response.SuccessResponse(c, http.StatusOK, "token refreshed", tokens)
}

// Sign out godoc
func (h *AuthHandler) SignOut(c *gin.Context) {
	log := logger.FromContext(c)
	log.Info("user signed out")
	response.SuccessResponse(c, http.StatusOK, "signed out successfully", nil)
}
