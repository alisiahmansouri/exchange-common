package auth_service

import (
	"context"
	"exchange-common/internal/entity"
	"exchange-common/internal/model"
	"exchange-common/internal/util"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	return s.repo.FindUserByPhone(ctx, phone)
}

func (s *Service) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.FindUserByEmail(ctx, email)
}

func (s *Service) FindUserByIdentifier(ctx context.Context, identifier string) (*entity.User, error) {
	if util.ValidateEmail(identifier) {
		return s.repo.FindUserByEmail(ctx, identifier)
	}
	if util.ValidatePhone(identifier) {
		return s.repo.FindUserByPhone(ctx, identifier)
	}
	return nil, model.ErrInvalidIdentifier
}

func (s *Service) Is2FAEnabled(ctx context.Context, userID uuid.UUID) (bool, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, model.ErrUserNotFound
	}
	return user.TwoFAEnabled, nil
}
