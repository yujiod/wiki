package recaptcha

import (
	"encoding/json"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Recaptcha struct {
	Success    string   `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// Recaptcha

func GetKey() (siteKey string, secretKey string) {
	recaptchaSiteKey := revel.Config.StringDefault("recaptcha.site_key", "")
	recaptchaSecretKey := revel.Config.StringDefault("recaptcha.secret_key", "")
	return recaptchaSiteKey, recaptchaSecretKey
}

func IsEnabled() bool {
	recaptchaSiteKey, recaptchaSecretKey := GetKey()
	return recaptchaSiteKey != "" && recaptchaSecretKey != ""
}

func IsAlways() bool {
	return revel.Config.BoolDefault("recaptcha.always", false)
}

func Validate(c *revel.Controller) bool {
	if IsEnabled() {
		recaptchaResponse := c.Params.Get("g-recaptcha-response")
		if recaptchaResponse == "" {
			return false
		}
		_, recaptchaSecretKey := GetKey()
		params := url.Values{}
		params.Add("secret", recaptchaSecretKey)
		params.Add("response", recaptchaResponse)
		params.Add("remoteip", c.Request.RemoteAddr)

		recaptchaUrl := "https://www.google.com/recaptcha/api/siteverify"
		response, err := http.PostForm(recaptchaUrl, params)
		defer response.Body.Close()
		recaptchaJson, _ := ioutil.ReadAll(response.Body)
		recaptchaResult := Recaptcha{}
		json.Unmarshal(recaptchaJson, &recaptchaResult)

		if err != nil || recaptchaResult.Success == "false" {
			return false
		}
	}
	return true
}
