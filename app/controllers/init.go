package controllers

import (
    "time"
    "github.com/revel/revel"
)

func init() {
    revel.OnAppStart(InitDB)
    revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
    revel.InterceptMethod((*GormController).Commit, revel.AFTER)
    revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)

    revel.TemplateFuncs["minus"] = func(a int, b int) int {
        return a - b
    }

    revel.TemplateFuncs["date"] = func(format string, time time.Time) string {
        //t, _ := time.Parse("2006-01-02 15:04:05.000000000 -0700 MST", current)
        return time.Local().Format(format)
    }
}