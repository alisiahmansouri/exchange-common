package auth_service

import (
	"context"
	"exchange-common/internal/entity"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	FindUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	UpdateUserEmailVerified(ctx context.Context, userID uuid.UUID, verified bool, verifiedAt *time.Time) error
	UpdateUserPhoneVerified(ctx context.Context, userID uuid.UUID, verified bool, verifiedAt *time.Time) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error
	DeleteUserByIDHard(ctx context.Context, userID uuid.UUID) error
}
