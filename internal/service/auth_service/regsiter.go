package auth_service

import (
	"context"
	"errors"
	"exchange-common/internal/entity"
	"exchange-common/internal/model"
	"exchange-common/internal/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser attempts to create a new user with either email or phone.
// It validates identifiers, checks for duplicates, hashes the password,
// and stores the user in the database. Returns the user ID if successful.
func (s *Service) RegisterUser(ctx context.Context, email, phone, fullName, password string) (*uuid.UUID, error) {
	// --- Normalize inputs ---
	email = util.NormalizeEmail(email)
	phone = util.NormalizePhone(phone)

	// --- Basic required check ---
	if email == "" && phone == "" {
		return nil, model.ErrEmailOrPhoneRequired
	}

	// --- Validate password strength ---
	if err := util.ValidatePassword(password); err != nil {
		return nil, err
	}

	// --- Check for existing email ---
	var emailPtr *string
	if email != "" {
		if !util.ValidateEmail(email) {
			return nil, model.ErrInvalidEmailFormat
		}
		existing, err := s.repo.FindUserByEmail(ctx, email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, model.ErrEmailExists
		}
		emailPtr = &email
	}

	// --- Check for existing phone ---
	var phonePtr *string
	if phone != "" {
		if !util.ValidatePhone(phone) {
			return nil, model.ErrInvalidPhone
		}
		existing, err := s.repo.FindUserByPhone(ctx, phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, model.ErrPhoneExists
		}
		phonePtr = &phone
	}

	// --- Hash password securely ---
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, model.ErrPasswordHash
	}

	// --- Create user entity ---
	user := &entity.User{
		Email:           emailPtr,
		Phone:           phonePtr,
		PasswordHash:    string(hashedPassword),
		FullName:        fullName,
		IsActive:        true,
		IsEmailVerified: false,
		IsPhoneVerified: false,
	}

	// --- Persist user in DB ---
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, model.ErrUserCreate
	}

	return &user.ID, nil
}
