package repository

import (
	"go-fiber-starter/utils"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func filterPaginate(modelName interface{}, pagination *utils.Pagination, clauses []clause.Expression) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	DB.Model(modelName).Clauses(clauses...).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(clauses...).Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

func setKeywordVarsByTotalExpr(keyword string, total int) (vars []interface{}) {
	for i := 0; i < total; i++ {
		vars = append(vars, strings.ToLower(keyword))
	}

	return vars
}

func lowerLikeQuery(field string) string {
	return "LOWER(" + field + ") LIKE '%'||?||'%'"
}
