package models

import (
	_ "github.com/revel/revel"
	"time"
)

// ページのリビジョン情報を持つ構造体
type Revision struct {
	Id           int
	PageId       int
	Title        string `sql:"size:255"`
	Body         string `sql:"size:16777215"`
	AddedLines   int
	DeletedLines int
	CreatedAt    time.Time
	DeletedAt    time.Time
}
