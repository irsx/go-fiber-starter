package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	GUID        uuid.UUID      `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SKU         string         `json:"sku" gorm:"type:varchar;uniqueIndex:idx_unique_product"`
	Name        string         `json:"name" gorm:"type:varchar"`
	Description string         `json:"description"`
	Image       string         `json:"image" gorm:"type:varchar"`
	Stock       int            `json:"stock" gorm:"type:integer;default:0"`
	SellPrice   float64        `json:"sell_price"`
	BuyPrice    float64        `json:"buy_price"`
	ExpiredAt   *string        `json:"expired_at" gorm:"type:date;default:null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (m *Product) TableName() string {
	return "product"
}

func (m *Product) GetImages() []string {
	var images = make([]string, 0)
	if m.Image != "" {
		json.Unmarshal([]byte(m.Image), &images)

		return images
	}

	images[0] = ""
	return images
}
