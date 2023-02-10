package transformer

import (
	"encoding/json"
	"go-fiber-starter/app/models"
	"go-fiber-starter/constants"
	"time"
)

type ProductListResponse struct {
	GUID      string  `json:"guid"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	UOM       string  `json:"uom"`
	Image     string  `json:"image"`
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
	Stock     int     `json:"stock"`
	ExpiredAt *string `json:"expired_at"`
}

func ProductListTransformer(Products []*models.Product) (rows []ProductListResponse) {
	for _, row := range Products {
		var mapResponse ProductListResponse
		jsonResponse, _ := json.Marshal(row)
		json.Unmarshal(jsonResponse, &mapResponse)

		if row.Image != "" {
			mapResponse.Image = row.GetImages()[0]
		}

		mapResponse.BuyPrice = row.BuyPrice
		mapResponse.SellPrice = row.SellPrice
		mapResponse.Stock = row.Stock

		if row.ExpiredAt != nil {
			expiredAt, _ := time.Parse(time.RFC3339, *row.ExpiredAt)
			expiredAtDate := expiredAt.Format(constants.DateFormat)
			mapResponse.ExpiredAt = &expiredAtDate
		}

		rows = append(rows, mapResponse)
	}

	return rows
}

type ProductDetailResponse struct {
	GUID        string    `json:"guid"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Stock       int       `json:"stock_in"`
	BuyPrice    float64   `json:"buy_price"`
	SellPrice   float64   `json:"sell_price"`
	ExpiredAt   *string   `json:"expired_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ProductDetailTransformer(product models.Product) (row ProductDetailResponse) {
	jsonResponse, _ := json.Marshal(product)
	json.Unmarshal(jsonResponse, &row)

	if row.ExpiredAt != nil {
		expiredAt, _ := time.Parse(time.RFC3339, *row.ExpiredAt)
		expiredAtDate := expiredAt.Format(constants.DateFormat)
		row.ExpiredAt = &expiredAtDate
	}

	return row
}
