package verification_service

import (
	"context"
	"exchange-common/internal/consts"
	"exchange-common/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) VerifyCode(ctx context.Context, userID uuid.UUID, channel, purpose, code string) error {
	if code == "" {
		return model.ErrVerificationCodeInvalid
	}
	if purpose == "" {
		return model.ErrVerificationPurposeInvalid
	}
	if channel != consts.ChannelEmail && channel != consts.ChannelSMS {
		return model.ErrVerificationChannelInvalid
	}
	return s.verifyGeneric(ctx, userID, code, purpose, channel)
}

func (s *Service) verifyGeneric(ctx context.Context, userID uuid.UUID, code, purpose, channel string) error {

	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return model.ErrInternal
	}

	vc, err := s.repo.FindVerificationByHashedCode(ctx, userID, string(hashedCode), purpose, channel)
	if err != nil {
		return model.ErrInternal
	}
	if vc == nil {
		return model.ErrVerificationCodeInvalidOrExpired
	}
	if vc.IsUsed {
		return model.ErrVerificationCodeAlreadyUsed
	}
	if vc.ExpiresAt.Before(time.Now()) {
		return model.ErrVerificationCodeExpired
	}

	vc.IsUsed = true
	if err := s.repo.UpdateVerificationCode(ctx, vc); err != nil {
		return model.ErrInternal
	}

	return nil
}
