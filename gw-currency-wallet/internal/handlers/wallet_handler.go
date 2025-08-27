package handlers

import (
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/pkg/logger"
	"gw-currency-wallet/internal/pkg/messages"
	"gw-currency-wallet/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	service *services.WalletService
}

func NewWalletHandler(service *services.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

// GetWallet godoc
// @Summary      Get wallet balances
// @Description  Get all currency balances for authenticated user
// @Tags         wallet
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} models.Response "Wallet balances retrieved"
// @Failure      401 {object} models.Response "Unauthorized"
// @Failure      500 {object} models.Response "Internal server error"
// @Router       /balance [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		logger.L.Error("Failed to get wallet")
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   messages.MsgUnauthorized,
		})

		return
	}

	wallet, err := h.service.GetWalletByUserID(c, userID)
	if err != nil {
		logger.L.Errorw("Failed to get wallet", "userID", userID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Error:   messages.MsgWalletNotFound,
			Details: err.Error(),
		})
		return
	}

	logger.L.Infow("Wallet retrieved", "userID", userID)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: wallet.GetAllBalances()})
}

// Deposit godoc
// @Summary      Deposit funds
// @Description  Deposit money to wallet in specified currency
// @Tags         wallet
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.WalletOperationReq true "Deposit data"
// @Success      200 {object} models.Response "Deposit successful"
// @Failure      400 {object} models.Response "Invalid request or deposit failed"
// @Failure      401 {object} models.Response "Unauthorized"
// @Router       /deposit [post]
func (h *WalletHandler) Deposit(c *gin.Context) {
	var req models.WalletOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Deposit request invalid", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgInvalidRequest,
			Details: err.Error(),
		})

		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   messages.MsgUnauthorized,
		})

		return
	}

	wallet, err := h.service.DepositWallet(c, userID, models.Currency(req.Currency), req.Amount)
	if err != nil {
		logger.L.Warnw("Deposit failed", "userID", userID, "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgDepositFailed,
			Details: err.Error(),
		})

		return
	}

	logger.L.Infow("Deposit successful", "userID", userID)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: wallet.GetAllBalances()})
}

// Withdraw godoc
// @Summary      Withdraw funds
// @Description  Withdraw money from wallet in specified currency
// @Tags         wallet
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.WalletOperationReq true "Withdraw data"
// @Success      200 {object} models.Response "Withdraw successful"
// @Failure      400 {object} models.Response "Invalid request or insufficient funds"
// @Failure      401 {object} models.Response "Unauthorized"
// @Failure      500 {object} models.Response "Internal server error"
// @Router       /withdraw [post]
func (h *WalletHandler) Withdraw(c *gin.Context) {
	var req models.WalletOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Withdraw request invalid", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgInvalidRequest,
			Details: err.Error(),
		})

		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   messages.MsgUnauthorized,
		})

		return
	}

	wallet, err := h.service.WithdrawWallet(c, userID, models.Currency(req.Currency), req.Amount)
	if err != nil {
		logger.L.Warnw("Withdraw failed", "userID", userID, "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgWithdrawFailed,
			Details: err.Error(),
		})

		return
	}

	logger.L.Infow("Withdraw successful", "userID", userID)
	c.JSON(http.StatusOK, models.Response{Success: true, Data: wallet.GetAllBalances()})
}

// Exchange godoc
// @Summary      Exchange currency
// @Description  Exchange money between currencies
// @Tags         wallet
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.ExchangeRequest true "Exchange data"
// @Success      200 {object} object "Exchange successful"
// @Failure      400 {object} models.Response "Invalid request or insufficient funds"
// @Failure      401 {object} models.Response "Unauthorized"
// @Router       /exchange [post]
func (h *WalletHandler) Exchange(c *gin.Context) {
	var req models.ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.L.Warnw("Exchange failed: invalid request", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgInvalidRequest,
			Details: err.Error(),
		})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Error:   messages.MsgUnauthorized,
		})
		return
	}

	wallet, err := h.service.ExchangeCurrency(c, userID, models.Currency(req.FromCurrency), models.Currency(req.ToCurrency), req.Amount)
	if err != nil {
		logger.L.Warnw("Exchange failed", "userID", userID, "from", req.FromCurrency, "to", req.ToCurrency, "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgExchangeFailed,
			Details: err.Error(),
		})
		return
	}

	balances := wallet.GetAllBalances()

	logger.L.Infow("Exchange successful", "userID", userID)
	c.JSON(http.StatusOK, gin.H{
		"message":          "Exchange successful",
		"exchanged_amount": req.Amount,
		"new_balance": gin.H{
			string(req.FromCurrency): balances[req.FromCurrency],
			string(req.ToCurrency):   balances[req.ToCurrency],
		},
	})
}

// GetAllRates godoc
// @Summary      Get exchange rates
// @Description  Get all available currency exchange rates
// @Tags         wallet
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} object "Rates retrieved"
// @Failure      401 {object} models.Response "Unauthorized"
// @Failure      500 {object} models.Response "Internal server error"
// @Router       /exchange/rates [get]
func (h *WalletHandler) GetAllRates(c *gin.Context) {
	mapRates, err := h.service.GetAllRates(c)
	if err != nil {
		logger.L.Warnw("Get rates failed", "error", err.Error())
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Error:   messages.MsgGetRatesFailed,
			Details: err.Error(),
		})
		return
	}

	logger.L.Infow("Get rates successful")
	c.JSON(http.StatusOK, gin.H{
		"rates": mapRates,
	})
}

func getUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		logger.L.Warn("user_id not found in context")
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.L.Warnw("Invalid user_id format", "userID", userIDStr, "error", err.Error())
		return uuid.Nil, false
	}

	return userID, true
}
