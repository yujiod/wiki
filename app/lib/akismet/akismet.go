package akismet

import (
	"fmt"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"net/url"
)

func IsEnabled() bool {
	return revel.Config.StringDefault("akismet.key", "") != ""
}

func Validate(c *revel.Controller, body string) bool {
	if IsEnabled() {
		akismetKey := revel.Config.StringDefault("akismet.key", "")
		params := url.Values{}
		params.Add("blog", "http://"+c.Request.Host)
		params.Add("user_ip", c.Request.RemoteAddr)
		params.Add("user_agent", c.Request.UserAgent())
		params.Add("referrer", c.Request.Referer())
		params.Add("comment_content", body)
		revel.INFO.Println(params.Encode())

		akismetUrl := fmt.Sprintf("https://%s.rest.akismet.com/1.1/comment-check", akismetKey)
		response, err := http.PostForm(akismetUrl, params)
		defer response.Body.Close()
		akismetResult, _ := ioutil.ReadAll(response.Body)
		revel.INFO.Println("Akismet result: %s", string(akismetResult))
		if err == nil && string(akismetResult) == "true" {
			return false
		}
	}
	return true
}
