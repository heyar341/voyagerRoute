package routes

import (
	"app/controllers/envhandler"
	"html/template"
	"net/http"
)

var  simulSearchTpl *template.Template

func init() {
	simulSearchTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/simul_search.html")).ParseGlob("templates/includes/*.html"))
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
