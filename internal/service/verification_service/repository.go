package verification_service

import (
	"context"
	"exchange-common/internal/entity"
	"github.com/google/uuid"
)

type Repository interface {
	FindVerificationByHashedCode(ctx context.Context, userID uuid.UUID, hashedCode, purpose, channel string) (*entity.VerificationCode, error)
	CreateVerification(ctx context.Context, code *entity.VerificationCode) error
	UpdateVerificationCode(ctx context.Context, code *entity.VerificationCode) error
	MarkVerificationCodeUsed(ctx context.Context, id uuid.UUID) error
	FindLatestActiveCode(ctx context.Context, userID uuid.UUID, identifier, purpose, channel string) (*entity.VerificationCode, error)
}
