package utils

import (
	"errors"
	"go-fiber-starter/constants"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func ExtractValidationError(req interface{}) error {
	var message string
	var v = validator.New()
	err := v.Struct(req)
	if err != nil {
		for i, err := range err.(validator.ValidationErrors) {
			if i > 0 {
				message += " | "
			}

			message += err.StructField() + ": " + err.Tag()
		}

		return errors.New(message)
	}

	return nil
}

func ValidateGUIDParams(ctx *fiber.Ctx) (string, error) {
	guid := ctx.Params("guid")
	if guid == "" {
		return "", JsonErrorValidation(ctx, constants.ErrGUIDRequired)
	}

	return guid, nil
}
