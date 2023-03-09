package dto

import (
	"database/sql"
	"go-fiber-starter/utils"
)

type ProductRequestDTO struct {
	UserGUID    string       // pass from jwt
	SKU         string       `json:"sku" validate:"required"`
	BarcodeID   string       `json:"barcode_id"`
	Name        string       `json:"name" validate:"required"`
	Description string       `json:"description" validate:"required"`
	ImageUrl    string       `json:"image,omitempty" validate:"required"`
	Stock       int          `json:"stock" validate:"numeric,min=0"`
	SellPrice   float64      `json:"sell_price" validate:"numeric,min=1"`
	BuyPrice    float64      `json:"buy_price" validate:"numeric,min=1"`
	ExpiredAt   sql.NullTime `json:"expired_at"`
}

func (req *ProductRequestDTO) Validate() error {
	return utils.ExtractValidationError(req)
}
