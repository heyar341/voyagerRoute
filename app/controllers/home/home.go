package home

import (
	"app/controllers"
	"encoding/base64"
	"html/template"
	"net/http"
)

type homeController struct {
	controllers.Controller
}

var homeTpl *template.Template

func init() {
	homeTpl = template.Must(template.Must(template.ParseGlob("templates/home/home.html")).ParseGlob("templates/includes/*.html"))
}
func showMsgWithCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		homeTpl.ExecuteTemplate(w, "home.html", data)
		return
	}
	data[c.Name] = string(b64Str)
	homeTpl.ExecuteTemplate(w, "home.html", data)
}

func Show(w http.ResponseWriter, req *http.Request) {
	var h homeController
	data := h.GetLoginStateFromCtx(req)
	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		showMsgWithCookie(w, c, data)
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		showMsgWithCookie(w, c, data)
		return
	}
	homeTpl.ExecuteTemplate(w, "home.html", data)
}
