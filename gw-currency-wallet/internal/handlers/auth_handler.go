package handlers

import (
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/pkg/logger"
	"gw-currency-wallet/internal/pkg/messages"
	"gw-currency-wallet/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService   *services.AuthService
	jwtManager    *services.JWTManager
	walletService *services.WalletService
}

func NewAuthHandler(authService *services.AuthService, walletService *services.WalletService, jwtManager *services.JWTManager) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		walletService: walletService,
		jwtManager:    jwtManager,
	}
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account with wallet
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Registration data"
// @Success      201 {object} models.Response "User created successfully"
// @Failure      400 {object} models.Response "Invalid request data"
// @Failure      409 {object} models.Response "User already exists or email already registered"
// @Failure      500 {object} models.Response "Internal server error"
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Invalid register request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgInvalidRequest,
			Details: err.Error(),
		})

		return
	}

	userID, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		logger.L.Warnw("Registration failed", "username", req.Username, "error", err.Error())

		var statusCode int
		var errorMsg string

		switch err {
		case services.ErrUserAlreadyExists:
			statusCode = http.StatusConflict
			errorMsg = messages.MsgUserExists
		case services.ErrEmailAlreadyExists:
			statusCode = http.StatusConflict
			errorMsg = messages.MsgEmailExists
		default:
			statusCode = http.StatusInternalServerError
			errorMsg = messages.MsgInternalError
		}

		c.JSON(statusCode, models.Response{
			Success: false,
			Error:   errorMsg,
		})

		return
	}

	_, err = h.walletService.CreateWallet(c.Request.Context(), userID)
	if err != nil {
		logger.L.Errorw("Failed to create wallet after registration",
			"userID", userID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   messages.MsgInternalError,
			Details: "User registered but wallet creation failed",
		})
		return
	}

	logger.L.Infow("User registered successfully", "userID", userID, "username", req.Username)
	c.JSON(http.StatusCreated, models.Response{
		Success: true,
		Data: gin.H{
			"user_id": userID.String(),
			"message": "User registered successfully",
		},
	})
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login credentials"
// @Success      200 {object} models.Response "Login successful"
// @Failure      400 {object} models.Response "Invalid request data"
// @Failure      401 {object} models.Response "Invalid credentials"
// @Failure      500 {object} models.Response "Internal server error"
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Invalid login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgInvalidRequest,
			Details: err.Error(),
		})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		logger.L.Warnw("Login failed", "username", req.Username, "error", err.Error())

		var statusCode int
		var errorMsg string

		switch err {
		case services.ErrInvalidCredentials:
			statusCode = http.StatusUnauthorized
			errorMsg = messages.MsgInvalidCredentials
		default:
			statusCode = http.StatusInternalServerError
			errorMsg = messages.MsgInternalError
		}

		c.JSON(statusCode, models.Response{
			Success: false,
			Error:   errorMsg,
		})
		return
	}

	logger.L.Infow("User logged successfully", "username", req.Username)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{
			"token": token,
		},
	})
}
