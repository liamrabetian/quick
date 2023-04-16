package service

import (
	"github.com/mohammadrabetian/quick/domain"
	"github.com/mohammadrabetian/quick/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *UserService) GetUserByToken(token string) (*domain.User, error) {
	return s.repo.GetUserByToken(token)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	return s.repo.UpdateUser(user)
}
