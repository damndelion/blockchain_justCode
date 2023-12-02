package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/damndelion/blockchain_justCode/internal/blockchain/metrics"
	"github.com/opentracing/opentracing-go"

	"github.com/damndelion/blockchain_justCode/pkg/cache"
	"github.com/skip2/go-qrcode"

	"github.com/damndelion/blockchain_justCode/config/blockchain"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/controller/http/middleware"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/controller/http/v1/dto"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/usecase"
	blockchainlogic "github.com/damndelion/blockchain_justCode/pkg/blockchain_logic"
	"github.com/damndelion/blockchain_justCode/pkg/logger"
	"github.com/gin-gonic/gin"
)

type chainRoutes struct {
	c          usecase.ChainUseCase
	l          logger.Interface
	cfg        *blockchain.Config
	chainCache cache.Blockchain
}

func newBlockchainRoutes(handler *gin.RouterGroup, c usecase.ChainUseCase, l logger.Interface, _ blockchainlogic.Blockchain, cfg *blockchain.Config, chainCache cache.Blockchain) {
	r := &chainRoutes{c, l, cfg, chainCache}

	blockchainHandler := handler.Group("/blockchain/wallet")
	{
		blockchainHandler.Use(middleware.JwtVerify(cfg.SecretKey))
		blockchainHandler.GET("/all", r.GetWallets)
		blockchainHandler.GET("/", r.GetWallet)
		blockchainHandler.GET("/balance", r.GetBalance)
		blockchainHandler.GET("/balance/address", r.GetBalanceByAddress)
		blockchainHandler.GET("/usd/balance", r.GetBalanceUSD)
		blockchainHandler.POST("/create", r.CreateWallet)
		blockchainHandler.POST("/send", r.Send)
		blockchainHandler.POST("/topup", r.TopUp)
		blockchainHandler.GET("/qr", r.GetWalletQRCode)
	}
}

// GetWallets godoc
// @Summary Get a list of wallets
// @Description Retrieve a list of wallets from the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {array} []string "List of wallets"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/all [get].
func (bc *chainRoutes) GetWallets(ctx *gin.Context) {
	span := opentracing.StartSpan("get wallets handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	wallets, err := bc.c.Wallets(spanCtx)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getWallets: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}

	ctx.JSON(http.StatusOK, wallets)
}

// GetWallet godoc
// @Summary Get a wallet by user ID
// @Description Retrieve a wallet from the blockchain for a specific user
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {string} string "Wallet address"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet [get].
func (bc *chainRoutes) GetWallet(ctx *gin.Context) {
	span := opentracing.StartSpan("get wallet handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")
	wallet, err := bc.c.Wallet(spanCtx, userID.(string))
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, "Wallet does not exist")

		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

// GetBalance godoc
// @Summary Get the balance of an address
// @Description Retrieve the balance of a specific address on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {number} float64 "Balance"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/balance [get].
func (bc *chainRoutes) GetBalance(ctx *gin.Context) {
	span := opentracing.StartSpan("get balance handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")
	balance, err := bc.c.GetBalance(spanCtx, userID.(string))
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalance: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}

	ctx.JSON(http.StatusOK, balance)
}

// GetBalanceByAddress godoc
// @Summary Get the balance of an address
// @Description Retrieve the balance of a specific address on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param address query string true "Wallet address"
// @Success 200 {number} float64 "Balance"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/balance/address [get].
func (bc *chainRoutes) GetBalanceByAddress(ctx *gin.Context) {
	span := opentracing.StartSpan("get balance by address handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	address := ctx.Query("address")
	balance, err := bc.c.GetBalanceByAddress(spanCtx, address)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalance: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}

	ctx.JSON(http.StatusOK, balance)
}

// GetBalanceUSD godoc
// @Summary Get the balance in USD of an address
// @Description Retrieve the balance in USD of a specific address on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {number} float64 "Balance in USD"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/usd/balance [get].
func (bc *chainRoutes) GetBalanceUSD(ctx *gin.Context) {
	span := opentracing.StartSpan("get balance is usd handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	balance, err := bc.c.GetBalanceUSD(spanCtx, userID.(string))
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - getBalanceUSD: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}

	ctx.JSON(http.StatusOK, balance)
}

// CreateWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {string} string "New Wallet"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/create [post].
func (bc *chainRoutes) CreateWallet(ctx *gin.Context) {
	span := opentracing.StartSpan("create wallet handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")

	wallet, err := bc.c.CreateWallet(spanCtx, userID.(string))
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - createWallet: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}
	ctx.JSON(http.StatusOK, wallet)
}

// Send godoc
// @Summary Send cryptocurrency to another address
// @Description Send cryptocurrency from one address to another on the blockchain
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param sendRequest body dto.SendRequest true "Send Request"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/send [post].
func (bc *chainRoutes) Send(ctx *gin.Context) {
	span := opentracing.StartSpan("send handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	var sendData dto.SendRequest
	err := ctx.ShouldBindJSON(&sendData)

	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("%v ", err))

		return
	}
	userID, _ := ctx.Get("user_id")

	err = bc.c.Send(spanCtx, userID.(string), sendData.To, sendData.Amount)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}
	metrics.TransactionRequestsTotalCollector.WithLabelValues(fmt.Sprintf("%v", ctx.Request.URL), strconv.Itoa(0), ctx.Request.Method).Inc()
	ctx.JSON(http.StatusOK, "Success ")
}

// TopUp godoc
// @Summary TopUp top up of an account
// @Description TopUp top up of an account
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Param topUpRequest body dto.TopupRequest true "Top up Request"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/topup [post].
func (bc *chainRoutes) TopUp(ctx *gin.Context) {
	span := opentracing.StartSpan("top up handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	var topupData dto.TopupRequest
	err := ctx.ShouldBindJSON(&topupData)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - topup: %w", err))
		errorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("%v ", err))

		return
	}
	userID, _ := ctx.Get("user_id")

	err = bc.c.TopUp(spanCtx, userID.(string), topupData.Amount)
	if err != nil {
		bc.l.Error(fmt.Errorf("http - v1 - blockchain - send: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}
	metrics.TransactionRequestsTotalCollector.WithLabelValues(fmt.Sprintf("%v", ctx.Request.URL), strconv.Itoa(0), ctx.Request.Method).Inc()
	ctx.JSON(http.StatusOK, "Success ")
}

// GetWalletQRCode godoc
// @Summary Get a wallet QR code by user ID
// @Description Return a wallet QR code by user ID
// @Tags Blockchain
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT Token"
// @Success 200 {string} string "Wallet QR code"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /v1/blockchain/wallet/qr [get].
func (bc *chainRoutes) GetWalletQRCode(ctx *gin.Context) {
	span := opentracing.StartSpan("get wallet qr code handler")
	defer span.Finish()
	spanCtx := opentracing.ContextWithSpan(ctx.Request.Context(), span)
	userID, _ := ctx.Get("user_id")
	wallet, err := bc.chainCache.Get(ctx, userID.(string))
	if err != nil {
		return
	}

	if wallet == "" {
		time.Sleep(1 * time.Second)
		wallet, err = bc.c.Wallet(spanCtx, userID.(string))
		if err != nil {
			bc.l.Error(fmt.Errorf("http - v1 - blockchain - GetWalletQRCode: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

			return
		}

		err = bc.chainCache.Set(ctx, userID.(string), wallet)
		if err != nil {
			bc.l.Error(fmt.Errorf("http - v1 - user - getUsersById: %w", err))
			errorResponse(ctx, http.StatusInternalServerError, "getUsersById cache error")
		}
	}

	qrCode, err := qrcode.Encode(wallet, qrcode.Medium, 256)
	if err != nil {
		bc.l.Error(fmt.Errorf("failed to generate QR code: %w", err))
		errorResponse(ctx, http.StatusInternalServerError, fmt.Sprintf("%v ", err))

		return
	}

	ctx.Header("Content-Type", "image/png")

	ctx.Data(http.StatusOK, "image/png", qrCode)
}
