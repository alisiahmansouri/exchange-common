package repository

import (
	"context"
	"errors"
	"exchange-common/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *GormRepo) CreateVerification(ctx context.Context, code *entity.VerificationCode) error {
	return getRepo(ctx, r).WithContext(ctx).Create(code).Error
}

func (r *GormRepo) UpdateVerificationCode(ctx context.Context, code *entity.VerificationCode) error {
	return getRepo(ctx, r).WithContext(ctx).Save(code).Error
}

func (r *GormRepo) MarkVerificationCodeUsed(ctx context.Context, id uuid.UUID) error {
	return getRepo(ctx, r).WithContext(ctx).
		Model(&entity.VerificationCode{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *GormRepo) FindVerificationByHashedCode(
	ctx context.Context,
	userID uuid.UUID,
	hashedCode, purpose, channel string,
) (*entity.VerificationCode, error) {
	var vc entity.VerificationCode

	err := getRepo(ctx, r).WithContext(ctx).
		Where("user_id = ? AND hashed_code = ? AND purpose = ? AND channel = ?",
			userID, hashedCode, purpose, channel).
		Order("created_at DESC").
		First(&vc).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vc, err
}

func (r *GormRepo) FindLatestActiveCode(ctx context.Context, userID uuid.UUID, identifier, purpose, channel string) (*entity.VerificationCode, error) {
	var vc entity.VerificationCode
	err := getRepo(ctx, r).WithContext(ctx).
		Where("user_id = ? AND identifier = ? AND purpose = ? AND channel = ? AND is_used = FALSE AND expires_at > NOW()",
			userID, identifier, purpose, channel).
		Order("created_at DESC").
		First(&vc).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vc, err
}
