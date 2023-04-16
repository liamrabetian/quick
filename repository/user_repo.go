package repository

import (
	"github.com/mohammadrabetian/quick/domain"
)

type UserRepository interface {
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByToken(token string) (*domain.User, error)
	UpdateUser(user *domain.User) error
}
