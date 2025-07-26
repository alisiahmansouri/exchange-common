package auth_service

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUserByIDHard(ctx, userID)
}
