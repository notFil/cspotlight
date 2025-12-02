package handlers

import (
	"net/http"

	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/pkg/auth"
	"github.com/notFil/cspotlight/pkg/errs"
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
		response.ErrorResponse(c, http.StatusBadRequest, "Authentication failed")
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
	response.Success(c, http.StatusOK, "login successful", tokens)
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
	var user models.UserRegisterDTO
	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid register payload", zap.Error(err))
		response.Error(c, err)
		return
	}

	if user.Password != user.ConfirmPassword {
		log.Warn("passwords do not match")
		response.Error(c, errs.ErrInvalidInput)
		return
	}

	log.Info("registering user", zap.String("email", user.Email))

	err := h.userService.RegisterUser(&user)
	if err != nil {
		log.Error("failed to register user", zap.String("email", user.Email), zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("user registered successfully", zap.String("email", user.Email))
	response.Success(c, http.StatusCreated, "user registered successfully", nil)
}

// Refresh token godoc
// @Summary      Refresh token
// @Description  Refreshes the JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RefreshTokenRequest  true  "Refresh token request"
// @Success      200   {object}  map[string]interface{} "token refreshed successfully"
// @Failure      401   {object}  map[string]interface{} "invalid refresh token"
// @Failure      500   {object}  map[string]interface{} "failed to refresh token"
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	log := logger.FromContext(c)
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BindJSON(&req); err != nil {
		log.Warn("invalid refresh token payload", zap.Error(err))
		response.Error(c, err)
		return
	}

	claims, err := auth.ValidateRefreshToken(req.RefreshToken, h.jwtConfig.RefreshSecretKey)
	if err != nil {
		log.Warn("invalid refresh token", zap.Error(err))
		response.Error(c, err)
		return
	}

	user, err := h.userService.GetUserByID(claims.Subject, *claims)
	if err != nil {
		log.Error("failed to get user for refresh token", zap.String("user_id", claims.Subject), zap.Error(err))
		response.Error(c, err)
		return
	}

	tokens, err := auth.GenerateJWT(user, h.jwtConfig)
	if err != nil {
		log.Error("failed to generate jwt for refresh", zap.String("user_id", user.ID), zap.Error(err))
		response.Error(c, err)
		return
	}

	log.Info("token refreshed", zap.String("user_id", user.ID))
	response.Success(c, http.StatusOK, "token refreshed", tokens)
}

// Sign out godoc
// @Summary      Sign out
// @Description  Signs out the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]interface{} "signed out successfully"
// @Failure      500   {object}  map[string]interface{} "failed to sign out"
// @Router       /api/auth/signout [post]
func (h *AuthHandler) SignOut(c *gin.Context) {
	log := logger.FromContext(c)
	log.Info("user signed out")
	response.Success(c, http.StatusOK, "signed out successfully", nil)
}
