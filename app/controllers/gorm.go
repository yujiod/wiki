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
	"regexp"
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
	dbSource := os.Getenv("DB_SOURCE")
	databaseUrl := ""

	if dbDriver == "" || dbSource == "" {
		if os.Getenv("DATABASE_URL") != "" {
			// Heroku Postgres
			databaseUrl = os.Getenv("DATABASE_URL")
		} else if os.Getenv("CLEARDB_DATABASE_URL") != "" {
			// Heroku ClearDB MySQL Database
			databaseUrl = os.Getenv("CLEARDB_DATABASE_URL")
		}

		if databaseUrl != "" {
			re, _ := regexp.Compile("([^:]+)://([^:]+):([^@]+)@([^/]+)/([^?]+)")
			match := re.FindStringSubmatch(databaseUrl)

			dbDriver = match[1]
			if dbDriver == "mysql" {
				dbSource = fmt.Sprintf(
					"%s:%s@tcp(%s:3306)/%s?parseTime=true",
					match[2],
					match[3],
					match[4],
					match[5],
				)
			} else {
				dbSource = databaseUrl
			}
		}
	}

	if dbDriver == "" {
		dbDriver = "sqlite3"
	}

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
