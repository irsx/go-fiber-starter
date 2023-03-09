package dto

import "go-fiber-starter/utils"

type UserRequestDTO struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (req *UserRequestDTO) Validate() error {
	return utils.ExtractValidationError(req)
}

type UserRegisterConsumerDTO struct {
	UserGUID    string `mapstructure:"userGUID"`
	Name        string `mapstructure:"name"`
	Email       string `mapstructure:"email"`
	Password    string `mapstructure:"password"`
	PhoneNumber string `mapstructure:"phoneNumber"`
}
