package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammadrabetian/quick/domain"
	"github.com/mohammadrabetian/quick/middleware"
)

var userSvc UserService

type UserService interface {
	GetUserByUsername(username string) (*domain.User, error)
	UpdateUser(user *domain.User) error
}

func InitUserHandlers(userService UserService) {
	userSvc = userService
}

type loginReqBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//	@Summary		Login API
//	@Description	Login authenticates the user and returns a token

//	@Tags		auth
//	@Accept		json
//	@Produce	json
//
//	@Param		_	body		loginReqBody	false	"username and password"
//
//	@Success	200	{object}	string
//	@Failure	400	{string}	httputil.HTTPError
//	@Router		/v1/auth/login [post]
func Login(c *gin.Context) {
	var req loginReqBody

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := userSvc.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if user == nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := middleware.GenerateSecureToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a token"})
		return
	}
	user.Token = token
	err = userSvc.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
