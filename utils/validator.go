package utils

import (
	"errors"
	"go-fiber-starter/constants"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func ExtractValidationError(req interface{}) error {
	var message string
	var v = validator.New()
	// get json tag
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	err := v.Struct(req)
	if err != nil {
		for i, err := range err.(validator.ValidationErrors) {
			if i > 0 {
				message += " | "
			}

			message += err.Field() + ": " + err.Tag()
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
