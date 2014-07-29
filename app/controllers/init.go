package controllers

import (
    "time"
    "github.com/revel/revel"
)

func init() {
    // 自動マイグレーション向けにサーバー起動時にInitDBを呼び出す
    revel.OnAppStart(InitDB)

    // 自動トランザクションを開始
    revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
    revel.InterceptMethod((*GormController).Commit, revel.AFTER)
    revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)

    // 四則演算を行うテンプレート関数を定義
    revel.TemplateFuncs["add"] = func(args ...int) int {
        result := 0
        for i, value := range args {
            if i == 0 {
                result = value
            } else {
                result += value
            }
        }
        return result
    }
    revel.TemplateFuncs["subtract"] = func(args ...int) int {
        result := 0
        for i, value := range args {
            if i == 0 {
                result = value
            } else {
                result -= value
            }
        }
        return result
    }
    revel.TemplateFuncs["multiply"] = func(args ...int) int {
        result := 0
        for i, value := range args {
            if i == 0 {
                result = value
            } else {
                result *= value
            }
        }
        return result
    }
    revel.TemplateFuncs["divide"] = func(args ...int) int {
        result := 0
        for i, value := range args {
            if i == 0 {
                result = value
            } else {
                result /= value
            }
        }
        return result
    }

    // 日付を書式化するテンプレート関数を定義
    // 書式化には 2006-01-02 15:04:05 -0700 MST の日時を指定する
    revel.TemplateFuncs["date"] = func(format string, time time.Time) string {
        return time.Local().Format(format)
    }
}