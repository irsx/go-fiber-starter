package controllers

import (
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type NewsController struct{}

func (c *NewsController) newsService() *services.NewsService {
	return new(services.NewsService)
}

func (c *NewsController) List(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ NEWS LIST")

	status, _ := strconv.Atoi(ctx.Query("status", "1"))
	return c.newsService().List(ctx, status)
}

func (c *NewsController) Add(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ NEWS ADD")

	req := new(dto.NewsRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	req.UserGUID = ctx.Locals("user_auth").(string)
	return c.newsService().Add(ctx, *req)
}

func (c *NewsController) Detail(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ NEWS DETAIL")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	return c.newsService().Detail(ctx, guid)
}

func (c *NewsController) Update(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ NEWS UPDATE")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	req := new(dto.NewsRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	return c.newsService().Update(ctx, guid, *req)
}

func (c *NewsController) Delete(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ NEWS DELETE")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	return c.newsService().Delete(ctx, guid)
}
