package repository

import (
	"errors"

	"github.com/mohammadrabetian/quick/domain"
	"gorm.io/gorm"
)

type userMySQLRepository struct {
	db *gorm.DB
}

func NewUserMySQLRepository(db *gorm.DB) UserRepository {
	return &userMySQLRepository{db: db}
}

func (r *userMySQLRepository) GetUserByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.Where("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userMySQLRepository) GetUserByToken(token string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.Where("token = ?", token).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userMySQLRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}
