package dto

import "go-fiber-starter/utils"

type NewsRequestDTO struct {
	UserGUID    string // pass from jwt
	Title       string `json:"title" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Description string `json:"description" validate:"required"`
	HyperLink   string `json:"hyperlink" validate:"required"`
	Status      int    `json:"status"`
}

func (req *NewsRequestDTO) Validate() error {
	return utils.ExtractValidationError(req)
}
