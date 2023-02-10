package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Page       int         `json:"page" query:"page"`
	Limit      int         `json:"limit" query:"limit"`
	SortBy     string      `json:"-" query:"sort_by"`
	SortDir    string      `json:"-" query:"sort_dir"`
	Keyword    string      `json:"-" query:"keyword"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"-"`
}

func GetPaginationParams(ctx *fiber.Ctx) Pagination {
	params := new(Pagination)
	if err := ctx.QueryParser(params); err != nil {
		return Pagination{}
	}

	return *params
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	sortDir := strings.ToUpper(p.SortDir)
	if p.SortBy == "" || (sortDir != "ASC" && sortDir != "DESC") {
		return ""
	}

	if sortDir == "" {
		p.SortDir = "ASC"
	}

	return p.SortBy + " " + sortDir
}
