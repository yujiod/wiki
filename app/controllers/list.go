package controllers

import (
	"github.com/revel/revel"
	"github.com/yujiod/wiki/app/models"
	"math"
	"strconv"
)

type List struct {
	App
}

// ページ一覧を表示する
func (c List) Index() revel.Result {
	order := c.Params.Get("order")
	if order == "" {
		order = "title"
	}

	paginateCurrent, _ := strconv.Atoi(c.Params.Get("page"))
	if paginateCurrent < 1 {
		paginateCurrent = 1
	}

	limit := 20

	pages := []models.Page{}

	db := c.db.Model(models.Page{})

	query := c.Params.Get("query")
	if query != "" {
		db = db.Where("title like ?", "%"+query+"%")
	}

	db.Order(order).Limit(limit).Offset(limit * (paginateCurrent - 1)).Find(&pages)

	paginateTotal := 0
	db.Count(&paginateTotal)

	paginateLast := int(math.Ceil(float64(paginateTotal) / float64(limit)))
	paginatePages := make([]int, paginateLast)

	return c.Render(pages, order, query, paginateCurrent, paginateTotal, paginateLast, paginatePages)
}
