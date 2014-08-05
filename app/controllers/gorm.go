package controllers

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"github.com/yujiod/wiki/app/models"
	"os"
)

type GormController struct {
	*revel.Controller
	db *gorm.DB
}

var (
	DB gorm.DB
)

// 自動マイグレーションを行う
func InitDB() {
	var err error

	dbDriver := os.Getenv("DB_DRIVER")
	if dbDriver == "" {
		dbDriver = "sqlite3"
	}

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "./wiki.db"
	}

	DB, err = gorm.Open(dbDriver, dbSource)

	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}

	DB.LogMode(revel.Config.BoolDefault("mode.dev", false))

	DB.AutoMigrate(models.Page{})
	DB.AutoMigrate(models.Revision{})

	DB.Model(models.Page{}).AddUniqueIndex("unique_title", "title")
}

// リクエスト時にトランザクションを開始する
func (c *GormController) Begin() revel.Result {
	c.db = DB.Begin()
	return nil
}

// リクエスト終了時にトランザクションを確定する
func (c *GormController) Commit() revel.Result {
	if c.db != nil {
		c.db.Commit()
	}
	c.db = nil
	return nil
}

// 異常時にトランザクションを破棄する
func (c *GormController) Rollback() revel.Result {
	if c.db != nil {
		c.db.Rollback()
	}
	return nil
}
