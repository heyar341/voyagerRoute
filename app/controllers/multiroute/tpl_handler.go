package multiroute

import (
	"app/controllers/envhandler"
	"app/model"
	"html/template"
	"net/http"
)

var multiSearchTpl, showRouteTpl *template.Template

func init() {
	multiSearchTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/search/multi_search.html")).ParseGlob("templates/includes/*.html"))
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/show_and_edit/multi_route_show.html")).ParseGlob("templates/includes/*.html"))
}

func MultiSearchTpl(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	data["apiKey"] = apiKey
	multiSearchTpl.ExecuteTemplate(w, "multi_search.html", data)
}

func ShowAndEditRoutesTpl(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	routeTitle := req.URL.Query().Get("route_title")
	user := req.Context().Value("user").(model.UserData)
	routeInfo := GetRoute(w, routeTitle, user.ID)
	data["apiKey"] = apiKey
	data["routeInfo"] = routeInfo
	showRouteTpl.ExecuteTemplate(w, "multi_route_show.html", data)
}
