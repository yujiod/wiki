package models

import (
	"github.com/revel/revel"
	"time"
)

// ページ情報を持つ構造体
type Page struct {
	Id        int
	Title     string `sql:"size:255"`
	Body      string `sql:"size:16777215"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func (page *Page) Validate(v *revel.Validation) {
	v.Check(
		page.Title,
		revel.Required{},
		revel.MaxSize{256},
	)
	v.Check(
		page.Body,
		revel.Required{},
	)
}
