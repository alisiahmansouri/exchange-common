package auth_service

import (
	"context"
	"errors"
	"exchange-common/internal/entity"
	"exchange-common/internal/model"
	"exchange-common/internal/util"

	"golang.org/x/crypto/bcrypt"
	"strings"
)

// AuthenticateUser verifies user credentials and returns the user.
func (s *Service) AuthenticateUser(ctx context.Context, identifier, password string) (*entity.User, error) {
	identifier = strings.TrimSpace(strings.ToLower(identifier))

	// اعتبارسنجی ورودی: ایمیل یا موبایل باید معتبر باشد
	if !(util.ValidateEmail(identifier) || util.ValidatePhone(identifier)) {
		return nil, errors.New("فرمت ایمیل یا شماره موبایل نامعتبر است")
	}

	if err := util.ValidatePassword(password); err != nil {
		return nil, err
	}

	var user *entity.User
	var err error

	if util.ValidateEmail(identifier) {
		user, err = s.repo.FindUserByEmail(ctx, identifier)
	} else {
		user, err = s.repo.FindUserByPhone(ctx, identifier)
	}
	if err != nil {
		return nil, model.ErrInvalidCreds
	}
	if user == nil {
		return nil, model.ErrInvalidCreds
	}

	if !user.IsActive {
		return nil, model.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, model.ErrInvalidCreds
	}

	return user, nil
}
