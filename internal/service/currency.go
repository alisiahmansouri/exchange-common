package service

import (
	"context"
	"exchange-common/internal/entity"
)

func (s *Service) ListActiveCurrencies(ctx context.Context) ([]entity.Currency, error) {
	return s.repo.ListActiveCurrencies(ctx)
}
