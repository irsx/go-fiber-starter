package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	GUID        uuid.UUID      `json:"guid" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserGUID    string         `json:"user_guid" gorm:"type:uuid"`
	SKU         string         `json:"sku" gorm:"uniqueIndex:product_ukey"`
	BarcodeID   string         `json:"barcode_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	Stock       int            `json:"stock"`
	SellPrice   float64        `json:"sell_price"`
	BuyPrice    float64        `json:"buy_price"`
	ExpiredAt   sql.NullTime   `json:"expired_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	User        User           `json:"-"`
}

func (m *Product) TableName() string {
	return "product"
}
