package home

import (
	"app/controllers"
	"app/internal/view"
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

func Show(w http.ResponseWriter, req *http.Request) {
	var h homeController
	data := h.GetLoginStateFromCtx(req)
	existsCookie := view.ExistsCookie(w, req, data, homeTpl, "home.html")
	if existsCookie {
		return
	}
	homeTpl.ExecuteTemplate(w, "home.html", data)
}
