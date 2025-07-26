package service

import (
	"context"
	"errors"
	"exchange-common/internal/entity"
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrDepositAmountInvalid  = errors.New("مبلغ واریز باید بزرگتر از صفر باشد")
	ErrWithdrawAmountInvalid = errors.New("مبلغ برداشت باید بزرگتر از صفر باشد")
	ErrWalletNotFound        = errors.New("کیف پول یافت نشد")
	ErrWalletUnauthorized    = errors.New("دسترسی به کیف پول غیرمجاز است")
	ErrInsufficientFunds     = errors.New("موجودی کافی نیست")
)

func (s *Service) ListWalletsByUser(ctx context.Context, userID uuid.UUID) ([]entity.Wallet, error) {
	return s.repo.ListWalletsByUserID(ctx, userID)
}

func (s *Service) Deposit(ctx context.Context, userID, currencyID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return ErrDepositAmountInvalid
	}

	err := s.repo.Transaction(ctx, func(txCtx context.Context) error {
		wallet, err := s.repo.FindWalletByUserIDAndCurrencyForUpdate(txCtx, userID, currencyID)
		if err != nil {
			return err
		}

		if wallet == nil {
			wallet = &entity.Wallet{
				UserID:     userID,
				CurrencyID: currencyID,
				Balance:    0,
				IsActive:   true,
			}
			if err := s.repo.CreateWallet(txCtx, wallet); err != nil {
				return err
			}
		}

		wallet.Balance += amount

		return s.repo.UpdateWallet(txCtx, wallet)
	})

	if err != nil {
		return fmt.Errorf("failed to deposit: %w", err)
	}

	return nil
}

func (s *Service) Withdraw(ctx context.Context, userID, walletID uuid.UUID, amount float64) error {
	if amount <= 0 {
		return ErrWithdrawAmountInvalid
	}

	err := s.repo.Transaction(ctx, func(txCtx context.Context) error {
		wallet, err := s.repo.FindWalletByIDForUpdate(txCtx, walletID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return ErrWalletNotFound
		}

		if wallet.UserID != userID {
			return ErrWalletUnauthorized
		}

		if wallet.Balance < amount {
			return ErrInsufficientFunds
		}

		wallet.Balance -= amount

		return s.repo.UpdateWallet(txCtx, wallet)
	})

	if err != nil {
		return fmt.Errorf("failed to withdraw: %w", err)
	}

	return nil
}
