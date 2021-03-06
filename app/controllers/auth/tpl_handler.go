package auth

import (
	"app/controllers"
	"encoding/base64"
	"html/template"
	"net/http"
)

var authTpl *template.Template

func init() {
	authTpl = template.Must(template.Must(template.ParseGlob("templates/auth/*.html")).ParseGlob("templates/includes/*.html"))
}

func processCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}, tplName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		authTpl.ExecuteTemplate(w, tplName, data)
		return
	}
	data[c.Name] = string(b64Str)
	authTpl.ExecuteTemplate(w, tplName, data)
}

func AskConfirmEmail(w http.ResponseWriter, req *http.Request) {
	data := controllers.GetLoginDataFromCtx(req)
	authTpl.ExecuteTemplate(w, "ask_confirm_email.html", data)
}

func RegisterForm(w http.ResponseWriter, req *http.Request) {
	data := controllers.GetLoginDataFromCtx(req)
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, data, "register.html")
		return
	}
	authTpl.ExecuteTemplate(w, "register.html", data)
}

func LoginForm(w http.ResponseWriter, req *http.Request) {
	data := controllers.GetLoginDataFromCtx(req)
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, data, "login.html")
		return
	}
	c, err = req.Cookie("success")
	if err == nil {
		processCookie(w, c, data, "login.html")
		return
	}
	authTpl.ExecuteTemplate(w, "login.html", data)
}
