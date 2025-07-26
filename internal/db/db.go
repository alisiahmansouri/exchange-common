package db

import (
	"context"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(gormDB *gorm.DB) *DB {
	return &DB{DB: gormDB}
}

func (d *DB) Ping(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
