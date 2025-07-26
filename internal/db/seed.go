package db

import (
	"exchange-common/internal/entity"
	"fmt"

	"gorm.io/gorm"
)

func SeedCurrencies(db *gorm.DB) error {
	currencies := []entity.Currency{
		{Code: "BTC", Name: "Bitcoin", Precision: 8, IsActive: true},
		{Code: "ETH", Name: "Ethereum", Precision: 18, IsActive: true},
		{Code: "USDT", Name: "Tether", Precision: 6, IsActive: true},
		{Code: "XRP", Name: "Ripple", Precision: 6, IsActive: true},
	}

	for _, c := range currencies {
		var count int64
		if err := db.Model(&entity.Currency{}).Where("code = ?", c.Code).Count(&count).Error; err != nil {
			return fmt.Errorf("error checking currency %s: %w", c.Code, err)
		}
		if count > 0 {
			continue // Already exists
		}
		if err := db.Create(&c).Error; err != nil {
			return fmt.Errorf("failed to create currency %s: %w", c.Code, err)
		}
	}

	return nil
}
