package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type defaultResponse struct {
	Success  bool        `json:"success"`
	Status   int         `json:"status"`
	Code     string      `json:"code,omitempty"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	Paginate *Pagination `json:"paginate,omitempty"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"description"`
	Message          string `json:"message"`
}

func JsonSuccess(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(defaultResponse{
		Success: true,
		Status:  fiber.StatusOK,
		Message: "OK",
		Data:    data,
	})
}

func JsonPagination(ctx *fiber.Ctx, pagination *Pagination) error {
	return ctx.Status(fiber.StatusOK).JSON(defaultResponse{
		Success:  true,
		Status:   fiber.StatusOK,
		Message:  "OK",
		Data:     pagination.Rows,
		Paginate: pagination,
	})
}

func JsonError(ctx *fiber.Ctx, err error, code string) error {
	errorMessage := logErrorFormat(err, code)
	Logger.Info(errorMessage)
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusBadRequest).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusBadRequest,
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	})
}

func JsonErrorInternal(ctx *fiber.Ctx, err error, code string) error {
	errorMessage := logErrorFormat(err, code)
	Logger.Info(errorMessage)
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusInternalServerError).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusInternalServerError,
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	})
}

func JsonErrorValidation(ctx *fiber.Ctx, err error) error {
	errorMessage := logErrorFormat(err, "E_VALIDATION")
	Logger.Info(errorMessage)
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusBadRequest).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusBadRequest,
		Code:    "E_VALIDATION",
		Message: err.Error(),
		Data:    nil,
	})
}

func JsonErrorNotFound(ctx *fiber.Ctx, err error) error {
	errorMessage := logErrorFormat(err, "E_NOT_FOUND")
	Logger.Info(errorMessage)
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusNotFound).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusNotFound,
		Code:    "E_NOT_FOUND",
		Message: err.Error(),
		Data:    nil,
	})
}

func JsonErrorUnauthorized(ctx *fiber.Ctx, err error) error {
	errorMessage := logErrorFormat(err, "E_UNAUTHORIZED")
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusUnauthorized).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusUnauthorized,
		Code:    "E_UNAUTHORIZED",
		Message: err.Error(),
		Data:    nil,
	})
}

func JsonErrorEnvironment(ctx *fiber.Ctx, env string) error {
	err := errors.New("missing env " + env + " variable")
	errorMessage := logErrorFormat(err, "E_ENV")
	Logger.Info(errorMessage)
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusInternalServerError).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusInternalServerError,
		Code:    "E_ENV",
		Message: err.Error(),
		Data:    nil,
	})
}

func JsonErrorForbidden(ctx *fiber.Ctx, err error) error {
	errorMessage := logErrorFormat(err, "E_FORBIDDEN")
	Logger.Info(errorMessage)
	Logger.Error(errorMessage)
	return ctx.Status(fiber.StatusForbidden).JSON(defaultResponse{
		Success: false,
		Status:  fiber.StatusForbidden,
		Code:    "E_FORBIDDEN",
		Message: err.Error(),
		Data:    nil,
	})
}

func logErrorFormat(err error, code string) string {
	return "❌ " + "[" + code + "] " + err.Error()
}
