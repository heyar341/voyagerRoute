package routes

import (
	"html/template"
	"net/http"
	"app/controllers/envhandler"
	"app/model"

)

var multiSearchTpl,simulSearchTpl,showRouteTpl *template.Template

func init()  {
	multiSearchTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/search/*")).ParseGlob("templates/includes/*.html"))
	simulSearchTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/*")).ParseGlob("templates/includes/*.html"))
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/show_and_edit/*")).ParseGlob("templates/includes/*.html"))
}


func MultiSearchTpl(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	data["apiKey"] = apiKey
	multiSearchTpl.ExecuteTemplate(w, "multi_search.html", data)
}

func SimulSearchTpl(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["apiKey"] = apiKey
	data["nineIterator"] = nineIterator
	simulSearchTpl.ExecuteTemplate(w, "simul_search.html", data)
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
