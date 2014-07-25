package controllers

import "github.com/revel/revel"

func init() {
    revel.OnAppStart(InitDB)
    revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
    revel.InterceptMethod((*GormController).Commit, revel.AFTER)
    revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)

    revel.TemplateFuncs["url"] = func(url string) string {
        return url
    }
}