package auth

import (
	"encoding/base64"
	"html/template"
	"net/http"
)

var authTpl *template.Template

func init() {
	authTpl = template.Must(template.Must(template.ParseGlob("templates/auth/*.html")).ParseGlob("templates/includes/*.html"))
}

func tplWithCookieMsg(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}, tplName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		authTpl.ExecuteTemplate(w, tplName, data)
		return
	}
	data[c.Name] = string(b64Str)
	authTpl.ExecuteTemplate(w, tplName, data)
}
