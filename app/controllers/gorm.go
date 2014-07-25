package controllers

import (
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "github.com/jinzhu/gorm"
    "github.com/revel/revel"
    "github.com/yujiod/wiki/app/models"
)

type GormController struct {
    *revel.Controller
    db gorm.DB
}

var (
    DB gorm.DB
)

func InitDB() {
    var err error
    DB, err = gorm.Open("sqlite3", "./app.db")

    if err != nil {
        panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
    }

    DB.LogMode(true)
    DB.AutoMigrate(models.Page{})

    DB.Model(models.Page{}).AddUniqueIndex("unique_title", "title")
}

func (c *GormController) Begin() revel.Result {
    c.db = DB
    return nil
}

func (c *GormController) Commit() revel.Result {
    c.db = DB
    return nil
}

func (c *GormController) Rollback() revel.Result {
    c.db = DB
    return nil
}
