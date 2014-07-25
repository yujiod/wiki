package controllers

import (
    "strconv"
    _ "regexp"
    _ "net/url"
    _ "html"
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

    // r, _ := regexp.Compile("\\[\\[((?:(?!\\]\\]).)+)\\]\\]")
    // body = r.ReplaceAllStringFunc(body, func(pageName string) string {
    //     return "<a href=\"/page/"+ url.QueryEscape(pageName) + "\">" + html.EscapeString(pageName) + "</a>"
    // })

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
    return c.Render(pageName, body)
}

func (c Page) Save(pageName string) revel.Result {
    page := models.Page{}
    c.db.Where("id = ?", pageName).Or("title = ?", pageName).First(&page)

    page.Title = pageName
    page.Body = c.Params.Get("PageBody")
    c.db.Save(&page)

    return c.Redirect("/page/"+pageName)
}
