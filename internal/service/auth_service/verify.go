package auth_service

import (
	"context"
	"time"

	"exchange-common/internal/model"
	"github.com/google/uuid"
)

// فرض: repo.UpdateUserEmailVerified(ctx, id, verified, verifiedAt) error در ریپازیتوری وجود دارد.
func (s *Service) MarkEmailVerified(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return model.ErrUserNotFound
	}
	if user.IsEmailVerified {
		return nil // اگر قبلاً تایید شده، خطا نده (idempotent)
	}
	verifiedAt := time.Now()
	return s.repo.UpdateUserEmailVerified(ctx, id, true, &verifiedAt)
}

// فرض: repo.UpdateUserPhoneVerified(ctx, userID uuid.UUID) error در ریپازیتوری وجود دارد.
func (s *Service) MarkPhoneVerified(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return model.ErrUserNotFound
	}
	if user.IsPhoneVerified {
		return nil // اگر قبلاً تایید شده، خطا نده (idempotent)
	}
	return s.repo.UpdateUserPhoneVerified(ctx, id)
}
