package models

import (
    "time"
    _ "github.com/revel/revel"
)

type Revision struct {
    Id   int
    PageId int
    Title string `sql:"size:255"`
    Body string `sql:"size:16777215"`
    AddedLines int
    DeletedLines int
    CreatedAt    time.Time
    DeletedAt    time.Time
}
