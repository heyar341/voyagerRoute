package simulsearch

import (
	"app/controllers/envhandler"
	"html/template"
	"log"
	"net/http"
)

var simulSearchTpl *template.Template

func init() {
	simulSearchTpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/simul_search.html")).ParseGlob("templates/includes/*.html"))
}

func SimulSearchTpl(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		log.Printf("Error whle gettibg data from context")
		msg := "エラーが発生しました。しばらく経ってからもう一度ご利用ください。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}
	//envファイルからAPIキー取得
	apiKey, err := envhandler.GetEnvVal("MAP_API_KEY")
	if err != nil {
		msg := "エラーが発生しました。しばらく経ってからもう一度ご利用ください。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}
	data["apiKey"] = apiKey
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["nineIterator"] = nineIterator
	simulSearchTpl.ExecuteTemplate(w, "simul_search.html", data)
}
