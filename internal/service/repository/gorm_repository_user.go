package repository

import (
	"exchange-common/internal/entity"
	"exchange-common/internal/model"
	"time"

	"errors"

	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ثبت کاربر جدید
func (r *GormRepo) CreateUser(ctx context.Context, user *entity.User) error {
	return getRepo(ctx, r).WithContext(ctx).Create(user).Error
}

// یافتن کاربر با ایمیل
func (r *GormRepo) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := getRepo(ctx, r).WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// یافتن کاربر با موبایل
func (r *GormRepo) FindUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := getRepo(ctx, r).WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// یافتن کاربر با آیدی
func (r *GormRepo) FindUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := getRepo(ctx, r).WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *GormRepo) UpdateUserPhoneVerified(ctx context.Context, userID uuid.UUID, verified bool, verifiedAt *time.Time) error {
	tx := getRepo(ctx, r).WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_phone_verified": verified,
			"phone_verified_at": verifiedAt,
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *GormRepo) UpdateUserEmailVerified(ctx context.Context, id uuid.UUID, verified bool, verifiedAt *time.Time) error {
	tx := getRepo(ctx, r).WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_email_verified": verified,
			"email_verified_at": verifiedAt,
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

// حذف کامل کاربر (hard delete)
func (r *GormRepo) DeleteUserByIDHard(ctx context.Context, userID uuid.UUID) error {
	return getRepo(ctx, r).WithContext(ctx).
		Unscoped().
		Where("id = ?", userID).
		Delete(&entity.User{}).Error
}

func (r *GormRepo) UpdateUserPassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	tx := getRepo(ctx, r).WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword) // فیلد صحیح طبق مدل
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
