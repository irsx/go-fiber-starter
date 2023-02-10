package controllers

import (
	"bufio"
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/sse"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
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
	utils.Logger.Info("user auth guid : " + userGUID)
	return c.ImportService().ImportProductExcel(ctx, userGUID)
}

func (c *ImportController) ImportProductsEvents(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ IMPORT PRODUCTS EVENTS")

	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")
	ctx.Set("Transfer-Encoding", "chunked")

	ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		sse.BrokerList.Subscribe(w)
		sse.BrokerList.Listen()
	}))

	return nil
}
