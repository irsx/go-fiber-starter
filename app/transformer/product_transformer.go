package transformer

import (
	"encoding/json"
	"go-fiber-starter/app/models"
	"go-fiber-starter/constants"
	"time"
)

type ProductListResponse struct {
	GUID      string `json:"guid"`
	BarcodeID string `json:"barcode_id"`
	SKU       string `json:"sku"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Stock     int    `json:"stock"`
}

func ProductListTransformer(Products []*models.Product) (rows []ProductListResponse) {
	for _, row := range Products {
		var mapResponse ProductListResponse
		jsonResponse, _ := json.Marshal(row)
		json.Unmarshal(jsonResponse, &mapResponse)
		rows = append(rows, mapResponse)
	}

	return rows
}

type ProductDetailResponse struct {
	GUID        string    `json:"guid"`
	BarcodeID   string    `json:"barcode_id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Stock       int       `json:"stock"`
	BuyPrice    float64   `json:"buy_price"`
	SellPrice   float64   `json:"sell_price"`
	ExpiredAt   string    `json:"expired_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ProductDetailTransformer(product models.Product) (row ProductDetailResponse) {
	jsonResponse, _ := json.Marshal(product)
	json.Unmarshal(jsonResponse, &row)

	if product.ExpiredAt.Valid {
		expiredAtDate := product.ExpiredAt.Time.Format(constants.DateFormat)
		row.ExpiredAt = expiredAtDate
	}

	return row
}
