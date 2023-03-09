package dto

import (
	"go-fiber-starter/utils"
)

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (req *LoginRequestDTO) Validate() error {
	return utils.ExtractValidationError(req)
}

type RegisterRequestDTO struct {
	Email       string `json:"email" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

func (req *RegisterRequestDTO) Validate() error {
	return utils.ExtractValidationError(req)
}
