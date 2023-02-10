package repository

import (
	"go-fiber-starter/app/models"
	"go-fiber-starter/app/transformer"
	"go-fiber-starter/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct{}

func (r *UserRepository) filterClauses(pagination *utils.Pagination) (clauses []clause.Expression) {
	if pagination.Keyword != "" {
		vars := setKeywordVarsByTotalExpr(pagination.Keyword, 2)
		query := lowerLikeQuery("name") + " OR " + lowerLikeQuery("national_id")
		clauses = append(clauses, clause.Expr{SQL: query, Vars: vars})
	}

	return clauses
}

func (r *UserRepository) GetAll(pagination utils.Pagination) (*utils.Pagination, error) {
	var users []*models.User
	clauses := r.filterClauses(&pagination)
	filter := filterPaginate(users, &pagination, clauses)
	if err := DB.Scopes(filter).Find(&users).Error; err != nil {
		return nil, err
	}

	pagination.Rows = transformer.UserListTransformer(users)
	return &pagination, nil
}

func (r *UserRepository) FindByNik(nik string) (user models.User, err error) {
	if err := DB.First(&user, "national_id = ?", nik).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (user models.User, err error) {
	if err := DB.First(&user, "email = ?", email).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) FindByGUID(uuid string) (user models.User, err error) {
	if err := DB.First(&user, "guid = ?", uuid).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) Insert(tx *gorm.DB, user models.User) (models.User, error) {
	if err := DB.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) UpdateByGUID(tx *gorm.DB, guid string, storeData models.User) (models.User, error) {
	var user models.User
	if err := tx.Model(&user).Clauses(clause.Returning{}).Where("guid = ?", guid).Updates(&storeData).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) UpdateByEmail(email string, storeData models.User) error {
	var user models.User
	if err := DB.Model(&user).Where("email = ?", email).Updates(&storeData).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteByGUID(guid string) error {
	var user models.User
	if err := DB.Where("guid = ?", guid).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
