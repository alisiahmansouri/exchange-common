package repository

import (
	"context"
	"errors"
	"exchange-common/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *GormRepo) CreateWallet(ctx context.Context, wallet *entity.Wallet) error {
	return getRepo(ctx, r).WithContext(ctx).Create(wallet).Error
}

func (r *GormRepo) FindWalletByUserIDAndCurrencyForUpdate(ctx context.Context, userID uuid.UUID, currencyID uuid.UUID) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := getRepo(ctx, r).WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND currency_id = ?", userID, currencyID).
		First(&wallet).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &wallet, err
}

func (r *GormRepo) FindWalletByIDForUpdate(ctx context.Context, id uuid.UUID) (*entity.Wallet, error) {
	var wallet entity.Wallet
	err := getRepo(ctx, r).WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &wallet, err
}

func (r *GormRepo) ListWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	err := getRepo(ctx, r).WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&wallets).Error
	return wallets, err
}

func (r *GormRepo) UpdateWallet(ctx context.Context, wallet *entity.Wallet) error {
	return getRepo(ctx, r).WithContext(ctx).Save(wallet).Error
}
