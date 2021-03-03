package multiroute

import (
	"app/controllers/envhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"log"
	"net/http"
)

var multiSearchTpl, showRouteTpl *template.Template

func init() {
	multiSearchTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/search/multi_search.html")).ParseGlob("templates/includes/*.html"))
	showRouteTpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/show_and_edit/multi_route_show.html")).ParseGlob("templates/includes/*.html"))
}

func MultiSearchTpl(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作を行ってください。"
	//envファイルからAPIキー取得
	apiKey, err := envhandler.GetEnvVal("MAP_API_KEY")
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}

	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	data["apiKey"] = apiKey
	multiSearchTpl.ExecuteTemplate(w, "multi_search.html", data)
}

func ShowAndEditRoutesTpl(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作を行ってください。"
	routeTitle := req.URL.Query().Get("route_title")
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	routeInfo, err := getRoute(routeTitle, user.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			msg = "ご指定いただいたルートがありません。"
			http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
			return
		} else {
		}
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		return
	}

	//envファイルからAPIキー取得
	apiKey, err := envhandler.GetEnvVal("MAP_API_KEY")
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	data["apiKey"] = apiKey
	data["routeInfo"] = routeInfo
	showRouteTpl.ExecuteTemplate(w, "multi_route_show.html", data)
}
