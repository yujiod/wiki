package controllers

import (
	"github.com/revel/revel"
	"github.com/yujiod/wiki/app/lib/wikihelper"
	"github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	// NewRelic agent
	agent *gorelic.Agent
)

func init() {
	// init for NewRelic
	initGorelic()
    revel.Filters = append(revel.Filters, gorelicFilter)

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

	revel.TemplateFuncs["len_gt"] = func(arg interface{}, length int) bool {
		if reflect.TypeOf(arg).Kind() == reflect.Slice {
			return reflect.ValueOf(arg).Len() > length
		}
		return false
	}
	revel.TemplateFuncs["len_lt"] = func(arg interface{}, length int) bool {
		if reflect.TypeOf(arg).Kind() == reflect.Slice {
			return reflect.ValueOf(arg).Len() < length
		}
		return false
	}
	revel.TemplateFuncs["replace"] = func(s string, old string, new string) string {
		return strings.Replace(s, old, new, -1)
	}
	revel.TemplateFuncs["urlencode"] = func(str string) string {
		return wikihelper.UrlEncode(str)
	}
}

// initGorelic Initializes the Gorelic agent
func initGorelic() {

	// Load NEWRELIC_LICENSE from a environment.
	// Only start gorelic if a license id resent.
	NEWRELIC_LICENSE := os.Getenv("NEWRELIC_LICENSE")
	if len(NEWRELIC_LICENSE) > 0 {
		log.Print("Starting newrelic daemon.")
		agent = gorelic.NewAgent()
		agent.NewrelicLicense = NEWRELIC_LICENSE
		agent.NewrelicName = "Wiki"
		agent.NewrelicPollInterval = 180
		agent.Verbose = true

		// "Manually" init the http timer (will be used in gorelicFilter)
		agent.CollectHTTPStat = true
		agent.HTTPTimer = metrics.NewTimer()

		agent.Run()
	} else {
		log.Print("!! Newrelic license missing from config file -> Not started")
	}
}

// Filter to capture HTTP metrics for gorelic
var gorelicFilter = func(c *revel.Controller, fc []revel.Filter) {
	startTime := time.Now()
	defer func() {
		if agent != nil {
			agent.HTTPTimer.UpdateSince(startTime)
		}
	}()
	fc[0](c, fc[1:])
}
