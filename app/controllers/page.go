package controllers

import (
    "strconv"
    "regexp"
    "net/url"
    "html"
    "crypto/sha1"
    "fmt"

    "github.com/revel/revel"
    "github.com/russross/blackfriday"
    "github.com/yujiod/wiki/app/models"
)

type Page struct {
    App
}

func (c Page) Show() revel.Result {
    var pageName string = c.Params.Get("pageName")
    if pageName == "" {
        pageName = "Home"
    }

    body := ""
    page := models.Page{}
    c.db.Where("id = ?", pageName).Or("title = ?", pageName).First(&page)

    id, _ := strconv.Atoi(pageName)

    if id != 0 && page.Title != "" {
        return c.Redirect("/page/"+page.Title)
    }

    if page.Body != "" {
        body = page.Body
    }

    re := regexp.MustCompile("\\[\\[([^\\]\\[]+)\\]\\]")
    body = re.ReplaceAllStringFunc(body, func(str string) string {
        str = regexp.MustCompile("(^\\[\\[|\\]\\]$)").ReplaceAllString(str, "")
        return "<a href=\"/page/"+ url.QueryEscape(str) + "\">" + html.EscapeString(str) + "</a>"
    })

    html := string(blackfriday.MarkdownCommon([]byte(body)))

    return c.Render(pageName, body, html)
}

func (c Page) Modify() revel.Result {
    var pageName string = c.Params.Get("pageName")

    body := ""
    page := models.Page{}
    c.db.Where("id = ?", pageName).Or("title = ?", pageName).First(&page)

    if page.Body != "" {
        body = page.Body
    }

    hash := fmt.Sprintf("%x", sha1.Sum([]byte(body)))

    return c.Render(pageName, body, hash)
}

func (c Page) Save(pageName string) revel.Result {
    page := models.Page{}
    c.db.Where("id = ?", pageName).Or("title = ?", pageName).First(&page)

    page.Title = pageName
    page.Body = c.Params.Get("PageBody")
    c.db.Save(&page)

    revision := models.Revision{}
    revision.Title = page.Title
    revision.Body = page.Body
    revision.PageId = page.Id
    c.db.Save(&revision)

    return c.Redirect("/page/"+pageName)
}

func (c Page) Revisions() revel.Result {
    var pageName string = c.Params.Get("pageName")

    page := models.Page{}
    c.db.Where("id = ?", pageName).Or("title = ?", pageName).First(&page)

    revisions := []models.Revision{}
    c.db.Where("page_id = ?", page.Id).Order("id desc").Find(&revisions)

    revisionSize := len(revisions)

    return c.Render(pageName, revisions, revisionSize)
}