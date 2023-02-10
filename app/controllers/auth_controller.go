package controllers

import (
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct{}

func (c *AuthController) authService() *services.AuthService {
	return new(services.AuthService)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ AUTH LOGIN")

	req := new(dto.LoginRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.authService().Authenticate(ctx, req)
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ AUTH REGISTER")

	req := new(dto.RegisterRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.authService().Register(ctx, *req)
}
