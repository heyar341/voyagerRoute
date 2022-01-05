package auth

import (
	"app/contexthandler"
	"encoding/base64"
	"html/template"
	"net/http"
)

var authTpl *template.Template

func init() {
	authTpl = template.Must(template.Must(template.ParseGlob("templates/auth/*.html")).ParseGlob("templates/includes/*.html"))
}

func processCookie(w http.ResponseWriter, c *http.Cookie, loginState map[string]interface{}, tplName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		authTpl.ExecuteTemplate(w, tplName, loginState)
		return
	}
	loginState[c.Name] = string(b64Str)
	authTpl.ExecuteTemplate(w, tplName, loginState)
}

func AskConfirmEmail(w http.ResponseWriter, req *http.Request) {
	loginState := contexthandler.GetLoginStateFromCtx(req)
	authTpl.ExecuteTemplate(w, "ask_confirm_email.html", loginState)
}

func RegisterForm(w http.ResponseWriter, req *http.Request) {
	loginState := contexthandler.GetLoginStateFromCtx(req)
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, loginState, "register.html")
		return
	}
	authTpl.ExecuteTemplate(w, "register.html", loginState)
}

func LoginForm(w http.ResponseWriter, req *http.Request) {
	loginState := contexthandler.GetLoginStateFromCtx(req)
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, loginState, "login.html")
		return
	}
	c, err = req.Cookie("success")
	if err == nil {
		processCookie(w, c, loginState, "login.html")
		return
	}
	authTpl.ExecuteTemplate(w, "login.html", loginState)
}
