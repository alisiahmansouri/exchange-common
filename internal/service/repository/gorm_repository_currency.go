package repository

import (
	"context"
	"errors"
	"exchange-common/internal/entity"
	"github.com/google/uuid"
	"strings"

	"gorm.io/gorm"
)

func (r *GormRepo) FindCurrencyByCode(ctx context.Context, code string) (*entity.Currency, error) {
	var currency entity.Currency
	err := getRepo(ctx, r).WithContext(ctx).
		Where("code = ?", strings.ToUpper(code)).
		First(&currency).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &currency, err
}

func (r *GormRepo) ListActiveCurrencies(ctx context.Context) ([]entity.Currency, error) {
	var currencies []entity.Currency
	err := getRepo(ctx, r).WithContext(ctx).
		Where("is_active = ?", true).
		Find(&currencies).Error
	return currencies, err
}

func (r *GormRepo) FindCurrencyByID(ctx context.Context, id uuid.UUID) (*entity.Currency, error) {
	var currency entity.Currency
	err := getRepo(ctx, r).WithContext(ctx).First(&currency, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &currency, err
}
