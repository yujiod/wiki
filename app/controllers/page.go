package controllers

import (
    "strconv"
    "regexp"
    "net/url"
    "html"
    "crypto/sha1"
    "fmt"
    "bytes"
    "strings"
    "github.com/pmezard/go-difflib/difflib"
    "github.com/revel/revel"
    "github.com/russross/blackfriday"
    "github.com/yujiod/wiki/app/models"
)

type Page struct {
    App
}

func (c Page) Show() revel.Result {
    pageName := c.Params.Get("pageName")

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

    revision := 0
    c.db.Model(models.Revision{}).Where("page_id = ?", page.Id).Count(&revision)

    return c.Render(pageName, body, html, page, revision)
}

func (c Page) Modify(pageName string) revel.Result {
    body := ""
    page := models.Page{}
    c.db.Where("title = ?", pageName).First(&page)

    if page.Body != "" {
        body = page.Body
    }

    hash := fmt.Sprintf("%x", sha1.Sum([]byte(body)))

    return c.Render(pageName, body, hash, page)
}

func (c Page) Save(pageName string) revel.Result {
    page := models.Page{}
    c.db.Where("title = ?", pageName).First(&page)

    body := c.Params.Get("page.Body")

    if page.Id > 0 && page.Body == body {
        return c.Redirect("/page/"+pageName)
    }

    page.Title = pageName
    page.Body = body
    c.db.Save(&page)

    previous := models.Revision{}
    c.db.Where("page_id = ?", page.Id).Order("id desc").First(&previous)

    unifiedDiff := difflib.UnifiedDiff{
        A:        difflib.SplitLines(html.EscapeString(previous.Body)),
        B:        difflib.SplitLines(html.EscapeString(page.Body)),
        Context:  65535,
    }
    diffString, _ := difflib.GetUnifiedDiffString(unifiedDiff)
    diffLines := difflib.SplitLines(diffString)

    revision := models.Revision{}

    for i, line := range diffLines {
        if i > 2 {
            if strings.HasPrefix(line, "+") {
                revision.AddedLines++;
            }
            if strings.HasPrefix(line, "-") {
                revision.DeletedLines++;
            }
        }
    }

    revision.Title = page.Title
    revision.Body = page.Body
    revision.PageId = page.Id
    c.db.Save(&revision)

    return c.Redirect("/page/"+pageName)
}

func (c Page) Revisions(pageName string) revel.Result {
    page := models.Page{}
    c.db.Where("title = ?", pageName).First(&page)

    revisions := []models.Revision{}
    c.db.Where("page_id = ?", page.Id).Order("id desc").Find(&revisions)

    revisionSize := len(revisions)

    return c.Render(pageName, revisions, revisionSize)
}

func (c Page) Diff(pageName string, revisionId string) revel.Result {
    page := models.Page{}
    c.db.Where("title = ?", pageName).First(&page)

    revision := models.Revision{}
    c.db.Where("page_id = ? and id = ?", page.Id, revisionId).First(&revision)

    previous := models.Revision{}
    c.db.Where("page_id = ? and id < ?", page.Id, revisionId).Order("id desc").First(&previous)

    unifiedDiff := difflib.UnifiedDiff{
        A:        difflib.SplitLines(html.EscapeString(previous.Body)),
        B:        difflib.SplitLines(html.EscapeString(revision.Body)),
        Context:  65535,
    }
    diffString, _ := difflib.GetUnifiedDiffString(unifiedDiff)
    diffLines := difflib.SplitLines(diffString)
    buffer := bytes.Buffer{}

    for i, line := range diffLines {
        if i > 2 {
            buffer.WriteString(line)
        }
    }

    diff := strings.TrimSpace(buffer.String())

    return c.Render(diff, revision, previous)
}
