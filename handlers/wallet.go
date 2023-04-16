package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohammadrabetian/quick/domain"
	"github.com/mohammadrabetian/quick/repository"
	"github.com/mohammadrabetian/quick/service"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var walletSvc WalletService

type WalletService interface {
	GetBalance(walletID uint64, username string) (*domain.Wallet, error)
	CreditWallet(walletID uint64, amount decimal.Decimal, username string) error
	DebitWallet(walletID uint64, amount decimal.Decimal, username string) error
}

func InitWalletHandlers(walletService WalletService) {
	walletSvc = walletService
}

//	@Summary		Get Balance API
//	@Description	Get Wallet Balance
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//
//	@Param			wallet_id	path		string	true	"wallet id to get the balance"
//
//	@Success		200			{object}	string
//	@Failure		400			{string}	httputil.HTTPError
//	@Router			/api/v1/wallets/{wallet_id}/balance [get]
//
//	@Security		ApiKeyAuth
func GetBalance(c *gin.Context) {
	walletID, err := strconv.ParseUint(c.Param("wallet_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}
	user := c.MustGet("user").(*domain.User)

	wallet, err := walletSvc.GetBalance(walletID, user.Username)
	if err != nil {
		logrus.Errorf("error in retrieving the wallet, err: %s", err)
		if errors.Is(err, repository.ErrWalletNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else if errors.Is(err, service.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Wallet does not belong to the user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve wallet balance"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": wallet.Balance})
}

type creditReqbody struct {
	Amount string `json:"amount"`
}

//	@Summary		Credit Wallet API
//	@Description	Creadit a Wallet Balance
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//
//	@Param			wallet_id	path		string			true	"wallet id to credit the balance"
//
//	@Param			_			body		creditReqbody	false	"amount to credit the wallet"
//
//	@Success		200			{object}	string
//	@Failure		400			{string}	httputil.HTTPError
//	@Router			/api/v1/wallets/{wallet_id}/credit [post]
//	@Security		ApiKeyAuth
func CreditWallet(c *gin.Context) {
	walletID, err := strconv.ParseUint(c.Param("wallet_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	req := creditReqbody{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-positive amount"})
		return
	}

	user := c.MustGet("user").(*domain.User)

	err = walletSvc.CreditWallet(walletID, amount, user.Username)
	if err != nil {
		logrus.Errorf("error in crediting the wallet, err: %s", err)
		if errors.Is(err, repository.ErrWalletNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else if errors.Is(err, service.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Wallet does not belong to the user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in crediting the wallet"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type debitReqbody struct {
	Amount string `json:"amount"`
}

//	@Summary		Debit Wallet API
//	@Description	Debit a Wallet Balance
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//
//	@Param			wallet_id	path		string			true	"wallet id to debit the balance"
//
//	@Param			_			body		debitReqbody	false	"amount to debit from the wallet"
//
//	@Success		200			{object}	string
//	@Failure		400			{string}	httputil.HTTPError
//	@Router			/api/v1/wallets/{wallet_id}/debit [post]
//	@Security		ApiKeyAuth
func DebitWallet(c *gin.Context) {
	walletID, err := strconv.ParseUint(c.Param("wallet_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
		return
	}

	req := debitReqbody{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil || amount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-positive amount"})
		return
	}

	user := c.MustGet("user").(*domain.User)

	err = walletSvc.DebitWallet(walletID, amount, user.Username)
	if err != nil {
		logrus.Errorf("error in debiting the wallet, err: %s", err)
		if errors.Is(err, repository.ErrWalletNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		} else if errors.Is(err, repository.ErrInsufficientFunds) {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient funds, balance cannot go below 0"})
		} else if errors.Is(err, service.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Wallet does not belong to the user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in debiting the wallet"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
