package domain

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	ID      uint64          `gorm:"primaryKey"`
	Balance decimal.Decimal `gorm:"type:decimal(64,8)"`
	UserID  string          `gorm:"index"`
}
