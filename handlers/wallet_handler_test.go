package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohammadrabetian/quick/domain"
	"github.com/mohammadrabetian/quick/handlers"
	"github.com/mohammadrabetian/quick/repository"
	"github.com/mohammadrabetian/quick/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetBalance(walletID uint64, username string) (*domain.Wallet, error) {
	args := m.Called(walletID, username)
	wallet, _ := args.Get(0).(*domain.Wallet)
	return wallet, args.Error(1)
}

func (m *MockWalletService) CreditWallet(walletID uint64, amount decimal.Decimal, username string) error {
	args := m.Called(walletID, amount, username)
	return args.Error(0)
}

func (m *MockWalletService) DebitWallet(walletID uint64, amount decimal.Decimal, username string) error {
	args := m.Called(walletID, amount, username)
	return args.Error(0)
}

func TestGetBalance(t *testing.T) {
	mockWalletSvc := new(MockWalletService)
	handlers.InitWalletHandlers(mockWalletSvc)
	gin.SetMode(gin.TestMode)

	t.Run("happy case", func(t *testing.T) {
		wallet := &domain.Wallet{ID: 1, Balance: decimal.NewFromInt(100)}
		mockWalletSvc.On("GetBalance", uint64(1), "user1").Return(wallet, nil)

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.GET("/wallets/:wallet_id/balance", handlers.GetBalance)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallets/1/balance", nil)
		req.Header.Set("user", "user1")
		r.ServeHTTP(w, req)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "100", response["balance"])
		mockWalletSvc.AssertExpectations(t)
	})

	t.Run("invalid wallet ID", func(t *testing.T) {
		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.GET("/wallets/:wallet_id/balance", handlers.GetBalance)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallets/invalid/balance", nil)
		req.Header.Set("user", "user1")
		r.ServeHTTP(w, req)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Invalid wallet ID", response["error"])
		mockWalletSvc.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockWalletSvc.On("GetBalance", uint64(2), "user1").Return(nil, repository.ErrWalletNotFound)

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})

		r.GET("/wallets/:wallet_id/balance", handlers.GetBalance)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallets/2/balance", nil)
		req.Header.Set("user", "user1")
		r.ServeHTTP(w, req)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "Wallet not found", response["error"])
		mockWalletSvc.AssertExpectations(t)
	})

	t.Run("wallet does not belong to user", func(t *testing.T) {
		mockWalletSvc.On("GetBalance", uint64(3), "user1").Return(nil, service.ErrAccessDenied)

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.GET("/wallets/:wallet_id/balance", handlers.GetBalance)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallets/3/balance", nil)
		req.Header.Set("user", "user1")
		r.ServeHTTP(w, req)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Equal(t, "Wallet does not belong to the user", response["error"])
		mockWalletSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockWalletSvc.On("GetBalance", uint64(4), "user1").Return(nil, errors.New("internal error"))

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.GET("/wallets/:wallet_id/balance", handlers.GetBalance)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/wallets/4/balance", nil)
		req.Header.Set("user", "user1")
		r.ServeHTTP(w, req)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Unable to retrieve wallet balance", response["error"])
		mockWalletSvc.AssertExpectations(t)
	})
}

func TestCreditWallet(t *testing.T) {

	mockWalletSvc := new(MockWalletService)
	handlers.InitWalletHandlers(mockWalletSvc)
	gin.SetMode(gin.TestMode)

	t.Run("successfully credit wallet", func(t *testing.T) {
		mockWalletSvc.On("CreditWallet", uint64(1), decimal.NewFromInt(100), "user1").Return(nil).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})

		r.PUT("/wallets/:wallet_id/credit", handlers.CreditWallet)

		w := httptest.NewRecorder()
		reqBody := `{"amount": "100"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/credit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockWalletSvc.AssertExpectations(t)
	})

	t.Run("invalid wallet ID", func(t *testing.T) {
		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.PUT("/wallets/:wallet_id/credit", handlers.CreditWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "100"}`
		req, _ := http.NewRequest("PUT", "/wallets/invalid_wallet_id/credit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid amount", func(t *testing.T) {
		mockWalletSvc.On("CreditWallet", uint64(1), decimal.NewFromInt(-100), "user1").Return(nil).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.PUT("/wallets/:wallet_id/credit", handlers.CreditWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "-100"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/credit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("wallet not found", func(t *testing.T) {
		mockWalletSvc.On("CreditWallet", uint64(1), decimal.NewFromInt(100), "user1").Return(repository.ErrWalletNotFound).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.PUT("/wallets/:wallet_id/credit", handlers.CreditWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "100"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/credit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("wallet does not belong to the user", func(t *testing.T) {
		mockWalletSvc.On("CreditWallet", uint64(1), decimal.NewFromInt(100), "user1").Return(service.ErrAccessDenied).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})
		r.PUT("/wallets/:wallet_id/credit", handlers.CreditWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "100"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/credit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

}

func TestDebitWallet(t *testing.T) {

	mockWalletSvc := new(MockWalletService)
	handlers.InitWalletHandlers(mockWalletSvc)
	gin.SetMode(gin.TestMode)

	t.Run("successfully debit wallet", func(t *testing.T) {
		mockWalletSvc.On("DebitWallet", uint64(1), decimal.NewFromInt(50), "user1").Return(nil).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})

		r.PUT("/wallets/:wallet_id/debit", handlers.DebitWallet)

		w := httptest.NewRecorder()
		reqBody := `{"amount": "50"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/debit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockWalletSvc.AssertExpectations(t)
	})

	t.Run("invalid wallet ID", func(t *testing.T) {
		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})

		r.PUT("/wallets/:wallet_id/debit", handlers.DebitWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "100"}`
		req, _ := http.NewRequest("PUT", "/wallets/invalid_wallet_id/debit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid amount", func(t *testing.T) {
		mockWalletSvc.On("DebitWallet", uint64(1), decimal.NewFromInt(-100), "user1").Return(nil).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})

		r.PUT("/wallets/:wallet_id/debit", handlers.DebitWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "-100"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/debit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		mockWalletSvc.On("DebitWallet", uint64(1), decimal.NewFromInt(100), "user1").Return(repository.ErrInsufficientFunds).Once()

		r := gin.Default()
		r.Use(func(c *gin.Context) {
			user := &domain.User{
				ID:       1,
				Username: "user1",
				Password: "password1",
			}
			c.Set("user", user)
			c.Next()
		})

		r.PUT("/wallets/:wallet_id/debit", handlers.DebitWallet)
		w := httptest.NewRecorder()
		reqBody := `{"amount": "100"}`
		req, _ := http.NewRequest("PUT", "/wallets/1/debit", strings.NewReader(reqBody))
		req.Header.Set("user", "user1")
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusPaymentRequired, w.Code)
	})
}
