package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/mohammadrabetian/quick/domain"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInsufficientFunds = errors.New("insufficient funds")

type walletMySQLRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

func (r *walletMySQLRepository) GetDB() *gorm.DB {
	return r.db
}
func NewWalletMySQLRepository(db *gorm.DB, cache *redis.Client) WalletRepository {
	return &walletMySQLRepository{db: db, cache: cache}
}

func (r *walletMySQLRepository) GetWallet(id uint64) (*domain.Wallet, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("wallet_%d", id)

	// Get the wallet object from cache
	walletJSON, err := r.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		wallet := &domain.Wallet{}
		if err := json.Unmarshal([]byte(walletJSON), wallet); err == nil {
			return wallet, nil
		}
	}

	wallet := &domain.Wallet{}

	err = r.db.First(wallet, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	// Cache the wallet object
	walletBytes, err := json.Marshal(wallet)
	if err == nil {
		r.cache.Set(ctx, cacheKey, string(walletBytes), 0)
	}

	return wallet, nil
}

func (r *walletMySQLRepository) UpdateWallet(tx *gorm.DB, id uint64, balance decimal.Decimal) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("wallet_%d", id)

	err := tx.Model(&domain.Wallet{}).Where("id = ?", id).Update("balance", balance).Error
	if err != nil {
		return err
	}

	// Invalidate the cache entry
	r.cache.Del(ctx, cacheKey)

	return nil
}
