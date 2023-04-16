package repository

import (
	"github.com/mohammadrabetian/quick/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWallet(id uint64) (*domain.Wallet, error)
	UpdateWallet(tx *gorm.DB, id uint64, balance decimal.Decimal) error
	GetDB() *gorm.DB
}
