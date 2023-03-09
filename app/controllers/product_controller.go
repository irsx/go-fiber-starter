package controllers

import (
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/services"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct{}

func (c *ProductController) productService() *services.ProductService {
	return new(services.ProductService)
}

func (c *ProductController) List(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ PRODUCT LIST")

	paginate := utils.GetPaginationParams(ctx)
	return c.productService().List(ctx, paginate)
}

func (c *ProductController) Detail(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ PRODUCT DETAIL")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	return c.productService().Detail(ctx, guid)
}

func (c *ProductController) Add(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ PRODUCT ADD")

	req := new(dto.ProductRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	req.UserGUID = ctx.Locals("user_auth").(string)
	return c.productService().Add(ctx, *req)
}

func (c *ProductController) Update(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ PRODUCT UPDATE")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	req := new(dto.ProductRequestDTO)
	if err := ctx.BodyParser(req); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	req.UserGUID = ctx.Locals("user_auth").(string)
	return c.productService().Update(ctx, guid, *req)
}

func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	utils.Logger.Info("✅ PRODUCT DELETE")

	guid, err := utils.ValidateGUIDParams(ctx)
	if err != nil {
		return err
	}

	return utils.JsonSuccess(ctx, fiber.Map{"guid": guid})
}
