package controllers

import (
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type UploadController struct{}

func (c *UploadController) uploadService() *services.UploadService {
	return new(services.UploadService)
}

func (c *UploadController) CDN(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ UPLOAD FILE TO CDN")

	fh, err := ctx.FormFile("file")
	if err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	imagePath, err := utils.UploadFileToStorage(fh)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_STORAGE_UPLOAD")
	}

	response, err := c.uploadService().UploadToCDN(ctx, imagePath)
	if err != nil {
		return err
	}

	return utils.JsonSuccess(ctx, &response)
}
