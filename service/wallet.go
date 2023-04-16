package service

import (
	"errors"

	"github.com/mohammadrabetian/quick/domain"
	"github.com/mohammadrabetian/quick/repository"
	"github.com/shopspring/decimal"
)

var ErrAccessDenied = errors.New("access denied")

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(walletID uint64, username string) (*domain.Wallet, error) {
	wallet, err := s.repo.GetWallet(walletID)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != username {
		return nil, ErrAccessDenied
	}

	return wallet, nil
}

func (s *WalletService) CreditWallet(walletID uint64, amount decimal.Decimal, username string) error {
	wallet, err := s.repo.GetWallet(walletID)
	if err != nil {
		return err
	}

	if wallet.UserID != username {
		return ErrAccessDenied
	}

	newBalance := wallet.Balance.Add(amount)

	tx := s.repo.GetDB().Begin()
	err = s.repo.UpdateWallet(tx, walletID, newBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *WalletService) DebitWallet(walletID uint64, amount decimal.Decimal, username string) error {
	wallet, err := s.repo.GetWallet(walletID)
	if err != nil {
		return err
	}

	if wallet.UserID != username {
		return ErrAccessDenied
	}

	newBalance := wallet.Balance.Sub(amount)
	if newBalance.IsNegative() {
		return repository.ErrInsufficientFunds
	}

	tx := s.repo.GetDB().Begin()
	err = s.repo.UpdateWallet(tx, walletID, newBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
