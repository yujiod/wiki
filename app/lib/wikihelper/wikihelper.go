package wikihelper

import (
	"github.com/russross/blackfriday"
	"net/url"
	"regexp"
	"strings"
)

func UrlEncode(title string) string {
	encodedTitle := strings.Replace(title, "/", "-", -1)
	encodedTitle = url.QueryEscape(encodedTitle)
	encodedTitle = strings.Replace(encodedTitle, "+", "%20", -1)
	return encodedTitle
}

func Render(markdown string) string {
	// ブラケットリンクを置換する
	body := markdown
	re := regexp.MustCompile("\\[\\[([^\\]\\[\\|]+)(\\|([^\\]\\[]+))?\\]\\]")
	for _, match := range re.FindAllStringSubmatch(body, -1) {
		bracketLink := match[0]
		title := match[1]
		alias := match[3]
		if alias == "" {
			alias = title
		}
		body = strings.Replace(body, bracketLink, "["+title+"](/page/"+UrlEncode(alias)+")", -1)
	}

	// Markdownへ変換
	html := string(blackfriday.MarkdownCommon([]byte(body)))
	html = strings.Replace(html, "<table>", "<table class=\"table table-bordered table-striped\">", -1)
	return html
}
