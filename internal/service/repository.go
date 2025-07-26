package service

import (
	"context"
	"exchange-common/internal/entity"
	"github.com/google/uuid"
)

type Repository interface {

	// Currency
	FindCurrencyByID(ctx context.Context, id uuid.UUID) (*entity.Currency, error)
	ListActiveCurrencies(ctx context.Context) ([]entity.Currency, error)

	// Wallet
	CreateWallet(ctx context.Context, wallet *entity.Wallet) error
	FindWalletByUserIDAndCurrencyForUpdate(ctx context.Context, userID uuid.UUID, currencyID uuid.UUID) (*entity.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *entity.Wallet) error
	ListWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Wallet, error)
	FindWalletByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Wallet, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
