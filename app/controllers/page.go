package controllers

import (
	"crypto/sha1"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/revel/revel"
	"github.com/yujiod/wiki/app/lib/wikihelper"
	"github.com/yujiod/wiki/app/models"
	"strconv"
	"strings"
)

type Page struct {
	App
}

// ページを表示する
func (c Page) Show() revel.Result {
	pageName := c.Params.Get("pageName")

	// ページ名が指定されていない場合はHomeとする
	if pageName == "" {
		pageName = "Home"
	}

	// ページ名もしくはページIDで検索
	page := models.Page{}
	c.db.Where("id = ?", pageName).Or("title = ?", pageName).First(&page)

	// ページIDが指定され、存在する場合はページ名のURLへリダイレクト
	id, _ := strconv.Atoi(pageName)
	if id != 0 && page.Id == id {
		return c.Redirect("/page/" + wikihelper.UrlEncode(page.Title))
	}

	// Markdownをレンダリング
	body := page.Body
	html := wikihelper.Render(body)

	// リビジョン番号を取得
	revision := 0
	c.db.Model(models.Revision{}).Where("page_id = ?", page.Id).Count(&revision)

	// 最近登録されたページ一覧を取得
	recentCreatedPages := []models.Page{}
	c.db.Order("created_at desc").Limit(10).Find(&recentCreatedPages)

	// 最近更新されたページ一覧を取得
	recentUpdatedPages := []models.Page{}
	c.db.Where("created_at != updated_at").Order("updated_at desc").Limit(10).Find(&recentUpdatedPages)

	return c.Render(pageName, body, html, page, revision, recentCreatedPages, recentUpdatedPages)
}

// ページを編集する
func (c Page) Modify(pageName string) revel.Result {
	// ページ名で検索
	page := models.Page{}
	c.db.Where("title = ?", pageName).First(&page)
	body := page.Body

	// 衝突検知のためのハッシュを生成
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(body)))

	return c.Render(pageName, body, hash, page)
}

// ページの登録もしくは更新を行う
func (c Page) Save(pageName string) revel.Result {
	// ページ名で検索
	page := models.Page{}
	c.db.Where("title = ?", pageName).First(&page)

	// POSTで送信された本文を取得
	body := c.Params.Get("page.Body")

	// ページは存在するが変更が一切ない場合には更新しない
	if page.Id > 0 && page.Body == body {
		return c.Redirect("/page/" + wikihelper.UrlEncode(page.Title))
	}

	// ページを保存する
	page.Title = pageName
	page.Body = body
	c.db.Save(&page)

	// 最新のリビジョンを取得
	previous := models.Revision{}
	c.db.Where("page_id = ?", page.Id).Order("id desc").First(&previous)

	// 追加行、削除行を数えるため差分を取得
	unifiedDiff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(previous.Body),
		B:       difflib.SplitLines(page.Body),
		Context: 65535,
	}
	diffString, _ := difflib.GetUnifiedDiffString(unifiedDiff)
	diffLines := difflib.SplitLines(diffString)

	// 追加行、削除行を数える
	revision := models.Revision{}
	for i, line := range diffLines {
		if i > 2 {
			if strings.HasPrefix(line, "+") {
				revision.AddedLines++
			}
			if strings.HasPrefix(line, "-") {
				revision.DeletedLines++
			}
		}
	}

	// リビジョンを保存
	revision.Title = page.Title
	revision.Body = page.Body
	revision.PageId = page.Id
	c.db.Save(&revision)

	return c.Redirect("/page/" + wikihelper.UrlEncode(pageName))
}

// ページのリビジョン一覧を表示する
func (c Page) Revisions(pageName string) revel.Result {
	// ページ名で検索
	page := models.Page{}
	c.db.Where("title = ?", pageName).First(&page)

	revisions := []models.Revision{}
	c.db.Where("page_id = ?", page.Id).Order("id desc").Find(&revisions)

	revisionSize := len(revisions)

	return c.Render(pageName, revisions, revisionSize)
}

// 指定したリビジョンIDとその直前のリビジョンの差分を表示する
// Ajaxでリクエストされることを前提とする
func (c Page) Diff(pageName string, revisionId string) revel.Result {
	// ページ名で検索
	page := models.Page{}
	c.db.Where("title = ?", pageName).First(&page)

	// 最新のリビジョンを取得
	revision := models.Revision{}
	c.db.Where("page_id = ? and id = ?", page.Id, revisionId).First(&revision)

	// 最新のリビジョンの直前のリビジョンを取得
	previous := models.Revision{}
	c.db.Where("page_id = ? and id < ?", page.Id, revisionId).Order("id desc").First(&previous)

	// 差分を生成
	unifiedDiff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(previous.Body),
		B:       difflib.SplitLines(revision.Body),
		Context: 65535,
	}
	diffString, _ := difflib.GetUnifiedDiffString(unifiedDiff)
	diffLines := difflib.SplitLines(diffString)

	// unified diff のヘッダーを除去する
	diffLines = diffLines[3:]

	// 編集前の内容が空の場合は、最初の行は空行を削除する差分なので削除
	if previous.Body == "" {
		diffLines = diffLines[1:]
	}
	diff := strings.Join(diffLines, "")
	diff = strings.TrimSpace(diff)

	// Markdownをレンダリング
	html := wikihelper.Render(revision.Body)

	return c.Render(diff, revision, previous, pageName, html)
}
