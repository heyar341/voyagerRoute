package auth

import (
	"html/template"
	"net/http"
)

var authTpl *template.Template

func init() {
	authTpl = template.Must(template.Must(template.ParseGlob("templates/auth/*.html")).ParseGlob("templates/includes/*.html"))
}

func AskConfirmEmail(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	authTpl.ExecuteTemplate(w, "ask_confirm_email.html", data)
}
func RegisterForm(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	data["qParams"] = req.URL.Query()
	authTpl.ExecuteTemplate(w, "register.html", data)
}
func LoginForm(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	data["qParams"] = req.URL.Query()
	authTpl.ExecuteTemplate(w, "login.html", data)
}
