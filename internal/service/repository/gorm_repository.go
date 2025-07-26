package repository

import (
	"context"

	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepo {
	return &GormRepo{db: db}
}

// Context key for storing/retrieving repo in ctx
type ctxRepoKeyType struct{}

var ctxRepoKey = ctxRepoKeyType{}

// Transaction starts a new transaction, sets the repo in ctx, and runs fn(ctx).
func (r *GormRepo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newRepo := &GormRepo{db: tx}
		ctx = context.WithValue(ctx, ctxRepoKey, newRepo)
		return fn(ctx)
	})
}

// RepoFromContext returns the transaction-aware repo from ctx (if present).
func RepoFromContext(ctx context.Context) (*GormRepo, bool) {
	repo, ok := ctx.Value(ctxRepoKey).(*GormRepo)
	return repo, ok
}

func (r *GormRepo) DB() *gorm.DB {
	return r.db
}

// Helper برای گرفتن db درست بر اساس ترنزکشن یا غیر ترنزکشن
func getRepo(ctx context.Context, r *GormRepo) *gorm.DB {
	if repo, ok := RepoFromContext(ctx); ok {
		return repo.DB()
	}
	return r.db
}
