package services

import (
	"errors"
	"go-fiber-starter/app/dto"
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/repository"
	"go-fiber-starter/app/transformer"
	"go-fiber-starter/constants"
	"go-fiber-starter/utils"

	"github.com/gofiber/fiber/v2"
)

type ProductService struct{}

func (s *ProductService) productRepo() *repository.ProductRepository {
	return new(repository.ProductRepository)
}

func (s *ProductService) List(ctx *fiber.Ctx, paginate utils.Pagination) error {
	paginationData, err := s.productRepo().GetAll(paginate)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_PRODUCT_LIST")
	}

	products := paginationData.Rows.([]*models.Product)
	paginationData.Rows = transformer.ProductListTransformer(products)
	return utils.JsonPagination(ctx, paginationData)
}

func (s *ProductService) Detail(ctx *fiber.Ctx, guid string) error {
	product, err := s.productRepo().FindByGUID(guid)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_PRODUCT_DETAIL")
	}

	row := transformer.ProductDetailTransformer(product)
	return utils.JsonSuccess(ctx, row)
}

func (s *ProductService) Add(ctx *fiber.Ctx, req dto.ProductRequestDTO) (err error) {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	productRepo := s.productRepo()
	if productRepo.IsExist(req.SKU, "") {
		return utils.JsonErrorValidation(ctx, errors.New("sku already exist"))
	}

	storeData := s.storeFromRequest(models.Product{}, req)
	product, err := productRepo.Insert(repository.DB, storeData)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_PRODUCT_ADD")
	}

	// SEND JOB
	utils.SendJobWithDefaultPayloads(constants.QueueNewProduct, product)

	return utils.JsonSuccess(ctx, product)
}

func (s *ProductService) Update(ctx *fiber.Ctx, GUID string, req dto.ProductRequestDTO) error {
	if err := req.Validate(); err != nil {
		return utils.JsonErrorValidation(ctx, err)
	}

	productRepo := s.productRepo()
	if productRepo.IsExist(req.SKU, GUID) {
		return utils.JsonErrorValidation(ctx, errors.New("sku already exist"))
	}

	product, err := productRepo.FindByGUID(GUID)
	if err != nil {
		return utils.JsonErrorNotFound(ctx, err)
	}

	storeData := s.storeFromRequest(product, req)
	updatedProduct, err := productRepo.UpdateByGUID(repository.DB, GUID, storeData)
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_PRODUCT_UPDATE")
	}

	updatedProduct.GUID = product.GUID
	// SEND JOB
	utils.SendJobWithDefaultPayloads(constants.QueueNewProduct, updatedProduct)

	return utils.JsonSuccess(ctx, updatedProduct)
}

func (s *ProductService) Delete(ctx *fiber.Ctx, guid string) error {
	if err := s.productRepo().DeleteByGUID(guid); err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_PRODUCT_DELETE")
	}

	return utils.JsonSuccess(ctx, fiber.Map{"guid": guid})
}

func (s *ProductService) storeFromRequest(Product models.Product, req dto.ProductRequestDTO) models.Product {
	Product.SKU = req.SKU
	Product.BarcodeID = req.BarcodeID
	Product.Name = req.Name
	Product.Description = req.Description
	Product.Image = req.ImageUrl
	Product.Stock = req.Stock
	Product.SellPrice = req.SellPrice
	Product.BuyPrice = req.BuyPrice
	Product.ExpiredAt = req.ExpiredAt
	return Product
}
