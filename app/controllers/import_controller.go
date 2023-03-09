package controllers

import (
	"go-fiber-starter/app/middlewares"
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type ImportController struct{}

func (c *ImportController) ImportService() *services.ImportService {
	return new(services.ImportService)
}

func (c *ImportController) ImportLogAll(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ IMPORT LOG ALL")

	return c.ImportService().GetLogs(ctx)
}

func (c *ImportController) ImportLogDetail(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ IMPORT LOG DETAIL")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	return c.ImportService().GetLogDetail(ctx, guid)
}

func (c *ImportController) ImportProducts(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ IMPORT PRODUCTS")

	userGUID := ctx.Locals("user_auth").(string)
	return c.ImportService().ImportProductExcel(ctx, userGUID)
}

func (c *ImportController) ImportProductsStream(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ IMPORT PRODUCTS STREAM")

	if err := middlewares.SseMiddleware(ctx); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_IMPORT_PRODUCTS_STREAM")
	}

	return nil
}
