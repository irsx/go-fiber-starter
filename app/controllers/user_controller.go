package controllers

import (
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (c *UserController) userService() *services.UserService {
	return new(services.UserService)
}

func (c *UserController) List(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ USER LIST")

	paginate := utils.GetPaginationParams(ctx)
	return c.userService().List(ctx, paginate)
}

func (c *UserController) Add(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ USER ADD")

	req := new(dto.UserRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.userService().Add(ctx, *req)
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ USER UPDATE")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	req := new(dto.UserRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.userService().Update(ctx, guid, *req)
}

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ USER DELETE")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	return c.userService().Delete(ctx, guid)
}
