package controllers

import (
	"database/sql"
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

	dbDriver := revel.Config.StringDefault("database.driver", "")
	dbSource := revel.Config.StringDefault("database.source", "")
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

	DB.AutoMigrate(&models.Page{}, &models.Revision{})

	DB.Model(models.Page{}).AddUniqueIndex("unique_title", "title")
}

// リクエスト時にトランザクションを開始する
func (c *GormController) Begin() revel.Result {
	if c.db != nil {
		return nil
	}
	db := DB.Begin()
	if db.Error != nil {
		panic(db.Error)
	}
	c.db = db
	return nil
}

// リクエスト終了時にトランザクションを確定する
func (c *GormController) Commit() revel.Result {
	if c.db == nil {
		return nil
	}
	c.db.Commit()
	if err := c.db.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.db = nil
	return nil
}

// 異常時にトランザクションを破棄する
func (c *GormController) Rollback() revel.Result {
	if c.db == nil {
		return nil
	}
	c.db.Rollback()
	if err := c.db.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.db = nil
	return nil
}
