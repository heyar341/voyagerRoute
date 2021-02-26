package mypage

import (
	"app/model"
	"html/template"
	"net/http"
)

var mypageTpl *template.Template

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

func ShowMypage(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	user := req.Context().Value("user").(model.UserData)
	data["userName"] = user.UserName
	mypageTpl.ExecuteTemplate(w, "mypage.html", data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	user := req.Context().Value("user").(model.UserData)
	titleNames := RouteTitles(user.ID)
	data["userName"] = user.UserName
	data["titles"] = titleNames
	mypageTpl.ExecuteTemplate(w, "show_routes.html", data)
}
