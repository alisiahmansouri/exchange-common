package verification_service

import (
	"context"
	"exchange-common/internal/consts"
	"exchange-common/internal/entity"
	"exchange-common/internal/model"
	"exchange-common/internal/util"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) sendVerificationCode(
	ctx context.Context,
	userID uuid.UUID,
	identifier, purpose, channel string,
	expireMinutes int,
	validateIdentifierFunc func(string) bool,
) error {
	if !validateIdentifierFunc(identifier) {
		return model.ErrVerificationIdentifierInvalid
	}

	activeCode, err := s.repo.FindLatestActiveCode(ctx, userID, identifier, purpose, channel)
	if err != nil {
		return model.ErrInternal
	}
	if activeCode != nil && activeCode.ExpiresAt.After(time.Now()) {
		// TODO: ارسال مجدد پیام (async/worker)
		return nil
	}

	code, err := util.GenerateVerificationCode()
	if err != nil {
		return model.ErrInternal
	}
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return model.ErrInternal
	}

	vc := &entity.VerificationCode{
		UserID:     userID,
		Identifier: identifier,
		HashedCode: string(hashedCode),
		ExpiresAt:  time.Now().Add(time.Duration(expireMinutes) * time.Minute),
		Purpose:    purpose,
		Channel:    channel,
		IsUsed:     false,
	}

	if err := s.repo.CreateVerification(ctx, vc); err != nil {
		return model.ErrInternal
	}

	// TODO: ارسال پیام (async) متناسب با کانال
	return nil
}

// Send2FACode sends a verification code to user's email or phone based on the channel.
// For SMS, it uses ValidatePhone; for email, it uses ValidateEmail.
func (s *Service) Send2FACode(ctx context.Context, userID uuid.UUID, identifier, channel, purpose string) error {
	switch channel {
	case consts.ChannelEmail:
		return s.sendVerificationCode(
			ctx, userID, identifier, purpose,
			consts.ChannelEmail, consts.Default2FAExpireMinutes, util.ValidateEmail,
		)
	case consts.ChannelSMS:
		return s.sendVerificationCode(
			ctx, userID, identifier, purpose,
			consts.ChannelSMS, consts.Default2FAExpireMinutes, util.ValidatePhone,
		)
	default:
		return model.ErrVerificationChannelInvalid
	}
}
