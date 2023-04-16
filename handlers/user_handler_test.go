package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mohammadrabetian/quick/domain"

	"github.com/gin-gonic/gin"
	"github.com/mohammadrabetian/quick/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	mockUserSvc := new(MockUserService)
	handlers.InitUserHandlers(mockUserSvc)
	gin.SetMode(gin.TestMode)

	t.Run("successful login", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "user1",
			Password: "password1",
			Token:    "",
		}
		mockUserSvc.On("GetUserByUsername", "user1").Return(user, nil).Once()
		mockUserSvc.On("UpdateUser", user).Return(nil).Once()

		r := gin.Default()
		r.POST("/login", handlers.Login)
		w := httptest.NewRecorder()
		reqBody := `{"username": "user1", "password": "password1"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserSvc.AssertExpectations(t)
	})

	t.Run("invalid username or password", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "user1",
			Password: "password1",
			Token:    "",
		}
		mockUserSvc.On("GetUserByUsername", "user1").Return(user, nil).Once()

		r := gin.Default()
		r.POST("/login", handlers.Login)
		w := httptest.NewRecorder()
		reqBody := `{"username": "user1", "password": "wrong_password"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("failed to retrieve user", func(t *testing.T) {
		mockUserSvc.On("GetUserByUsername", "user1").Return(nil, errors.New("database error")).Once()

		r := gin.Default()
		r.POST("/login", handlers.Login)
		w := httptest.NewRecorder()
		reqBody := `{"username": "user1", "password": "password1"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("failed to update user token", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "user1",
			Password: "password1",
			Token:    "",
		}
		mockUserSvc.On("GetUserByUsername", "user1").Return(user, nil).Once()
		mockUserSvc.On("UpdateUser", user).Return(errors.New("database error")).Once()

		r := gin.Default()
		r.POST("/login", handlers.Login)
		w := httptest.NewRecorder()
		reqBody := `{"username": "user1", "password": "password1"}`
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
