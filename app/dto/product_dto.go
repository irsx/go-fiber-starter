package dto

import "go-fiber-starter/utils"

type ProductRequestDTO struct {
	CreatedBy   string
	SKU         string  `json:"sku" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	ImageUrl    string  `json:"image,omitempty" validate:"required"`
	Stock       int     `json:"stock_in" validate:"numeric,min=0"`
	SellPrice   float64 `json:"sell_price" validate:"numeric,min=1"`
	BuyPrice    float64 `json:"buy_price" validate:"numeric,min=1"`
	ExpiredAt   *string `json:"expired_at"`
}

func (req *ProductRequestDTO) Validate() error {
	return utils.ExtractValidationError(req)
}
