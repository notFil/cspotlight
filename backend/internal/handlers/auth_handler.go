package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/models"
	"github.com/notFil/cspotlight/internal/response"
	"github.com/notFil/cspotlight/internal/services"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
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
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var authRequest models.AuthRequest
	if err := c.BindJSON(&authRequest); err != nil {
		log.Warn("invalid login payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid login payload"))
		return
	}

	log.Info("attempting login", zap.String("username", authRequest.Username))

	user, err := h.userService.AuthenticateUser(ctx, &authRequest)
	if err != nil {
		log.Warn("authentication failed", zap.String("username", authRequest.Username), zap.Error(err))
		c.Error(apperrors.New(http.StatusUnauthorized, "authentication failed"))
		return
	}

	session := sessions.Default(c)
	session.Set("userID", user.ID.String())
	session.Set("role", user.Role)
	session.Set("teamID", user.TeamID.String())
	if err = session.Save(); err != nil {
		log.Error("failed to save session", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "authentication failed"))
		return
	}

	log.Info("login successful", zap.String("user_id", user.ID.String()))
	response.Success(c, http.StatusOK, "login successful", user)
}

// Register godoc
// @Summary      Register user
// @Description  Registers a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.UserRegisterDTO  true  "User registration info"
// @Success      201   {object}  map[string]interface{} "User registered successfully"
// @Failure      400   {object}  map[string]interface{} "Invalid request payload"
// @Failure      500   {object}  map[string]interface{} "Failed to register user"
// @Router       /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	var user models.UserRegisterDTO
	if err := c.BindJSON(&user); err != nil {
		log.Warn("invalid register payload", zap.Error(err))
		c.Error(apperrors.New(http.StatusBadRequest, "invalid registration payload"))
		return
	}

	log.Info("registering user", zap.String("email", user.Email))

	if err := h.userService.RegisterUser(ctx, &user); err != nil {
		log.Error("failed to register user", zap.String("email", user.Email), zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to register user"))
		return
	}

	log.Info("user registered successfully", zap.String("email", user.Email))
	response.Success(c, http.StatusCreated, "user registered successfully", nil)
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
	log.Info("user attempting to sign out")

	ctx := c.Request.Context()
	user := auth.GetUserContext(ctx)
	if user == nil {
		log.Warn("no user found in context")
		c.Error(apperrors.New(http.StatusUnauthorized, "no user found in context"))
		return
	}

	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: -1})
	session.Clear()
	if err := session.Save(); err != nil {
		log.Error("failed to save session", zap.Error(err))
		c.Error(apperrors.New(http.StatusInternalServerError, "failed to sign out"))
		return
	}

	response.Success(c, http.StatusOK, "signed out successfully", nil)
}
