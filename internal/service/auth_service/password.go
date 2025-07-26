package auth_service

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) ResetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserPassword(ctx, userID, string(hashed))
}
