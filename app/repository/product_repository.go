package repository

import (
	"go-fiber-starter/app/models"
	"go-fiber-starter/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct{}

func (r *ProductRepository) filterClauses(keyword string) (clauses []clause.Expression) {
	if keyword != "" {
		vars := setKeywordVarsByTotalExpr(keyword, 2)
		query := lowerLikeQuery("name") + " OR " + lowerLikeQuery("sku")
		clauses = append(clauses, clause.Expr{SQL: query, Vars: vars, WithoutParentheses: true})
	}

	return clauses
}

func (r *ProductRepository) GetAll(pagination utils.Pagination) (*utils.Pagination, error) {
	var Products []*models.Product
	clauses := r.filterClauses(pagination.Keyword)
	filter := filterPaginate(Products, &pagination, clauses)
	if err := DB.Scopes(filter).Find(&Products).Error; err != nil {
		return nil, err
	}

	pagination.Rows = Products
	return &pagination, nil
}

func (r *ProductRepository) FindByGUID(guid string) (Product models.Product, err error) {
	if err := DB.First(&Product, "guid = ?", guid).Error; err != nil {
		return Product, err
	}

	return Product, nil
}

func (r *ProductRepository) Insert(tx *gorm.DB, Product models.Product) (models.Product, error) {
	if err := tx.Create(&Product).Error; err != nil {
		return Product, err
	}

	return Product, nil
}

func (r *ProductRepository) InsertMany(tx *gorm.DB, products []models.Product, batchSize int) error {
	if err := DB.CreateInBatches(&products, batchSize).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) UpdateByGUID(tx *gorm.DB, guid string, storeData models.Product) (models.Product, error) {
	var Product models.Product
	if err := tx.Model(&Product).Where("guid = ?", guid).Updates(&storeData).Error; err != nil {
		return Product, err
	}

	return Product, nil
}

func (r *ProductRepository) DeleteByGUID(guid string) error {
	var Product models.Product
	if err := DB.Where("guid = ?", guid).Delete(&Product).Error; err != nil {
		return err
	}

	return nil
}

func (s *ProductRepository) IsExist(SKU string, GUID string) bool {
	var product models.Product
	if err := DB.Unscoped().Select("guid").First(&product, "sku = ?", SKU).Error; err != nil {
		return false
	}

	if GUID != "" && product.GUID.String() == GUID {
		return false
	}

	return true
}
